package server

import (
	"errors"
	"github.com/axengine/cache"
	"github.com/axengine/cache/persistence"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/ginprom"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wolot/api/config"
	_ "github.com/wolot/api/docs"
	"github.com/wolot/api/libs"
	"github.com/wolot/api/libs/ginlimiter"
	"github.com/wolot/api/web/dbo"
	"github.com/wolot/api/web/handlers"
	"github.com/wolot/api/web/proxy"
	"github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
	"os"
	"time"
)

type Server struct {
	cfg     *config.Config
	handler *handlers.Handler
	proxy   *proxy.ReverseProxy
	metrics *ginprom.GinPrometheus
}

func NewServer(logger *zap.Logger, cfg *config.Config, dbo3 *dbo.DBO) *Server {
	handler := handlers.NewHandler(logger, cfg, dbo3)

	p := proxy.NewReverseProxy()
	p.AddToSetUpstream(cfg.RPC)
	//  '/' tendermint 会把所有路由打印出来，会暴露源endpoint，无法拦截请求。
	p.AddToSetLimitPath("/", "/broadcast_tx_async", "/broadcast_tx_sync", "/broadcast_tx_commit")

	return &Server{
		cfg:     cfg,
		handler: handler,
		proxy:   p,
	}
}

// @title MONDO API
// @version 3.0.0
// @description MONDO 区块链API文档
// @host https://oloapi-test.wolot.io
// @BasePath /
func (s *Server) Start() {
	// 生产环境可以使用 gin.New()
	router := gin.New()

	// pprof
	pprof.Register(router)

	// 日志打印中间件
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/", "/metrics", "/health"}}))

	// cross domain
	router.Use(libs.Cors())

	router.Any("/", func(context *gin.Context) {
		context.String(200, "%s", "ok")
	})

	// 开启metrics
	if s.cfg.Metrics {
		p := ginprometheus.NewPrometheus("gin")
		p.Use(router)
	}
	router.Static("/static", "./static")
	if s.cfg.Dev {
		router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	store := persistence.NewInMemoryStore(time.Second)

	// v2版普通接口
	v2 := router.Group("/v2")
	{
		v2.GET("/genkey", s.handler.V2GenKey)                      //生成帐户，不上链
		v2.GET("/accounts/:address", s.handler.V2QueryAccount)     //根据地址查询帐户
		v2.GET("/convert/:publickey", s.handler.V2Convert)         //公钥生成地址
		v2.POST("/transactions", s.handler.SignedBatchTransaction) //发起批量交易(1vN)
	}

	lm := ginlimiter.NewRateLimiter(s.cfg.Limiter.Interval.Duration, s.cfg.Limiter.Capacity, func(ctx *gin.Context) (string, error) {
		//key := ctx.Request.Header.Get("X-API-KEY")
		key := ctx.Request.RequestURI
		if key != "" {
			return key, nil
		}
		return "", errors.New("key is missing")
	})

	contract := router.Group("/v2/contract")
	{
		contract.GET("/accounts/:address", s.handler.QueryContract)                                         //查询合约帐户
		contract.GET("/events/:txhash", cache.CachePageAtomic(store, time.Minute, s.handler.QueryTxEvents)) //查询tx events(events) from statedb
		contract.POST("/transactions", s.handler.SignedEvmTransaction)                                      //发送合约签名交易(创建/执行/call,call要消耗gas)
		contract.POST("/sign", s.handler.SignEvmTransaction)                                                // 签名接口
		contract.POST("/query", lm.Middleware(), s.handler.ContractSignedCallTx)                            //签名query合约(evm本地执行，不消耗gas，不上链)
		contract.POST("/multisigTransactions", s.handler.SignedEvmMutlisigTransaction)                      // 多签交易
		contract.POST("/multisigner", s.handler.Multisigner)                                                // 取多签签名者

		contract.POST("/deploy", s.handler.ContractDeployTx) //部署合约
		contract.POST("/invoke", s.handler.ContractInvokeTx) //调用合约
		contract.POST("/call", s.handler.ContractCallTx)     //call合约(evm本地执行，不消耗gas，不上链)
	}

	erc20Group := router.Group("/v2/erc20")
	{
		erc20Group.GET("/:token/balanceOf/:to", s.handler.BalanceOf)
	}

	v3 := router.Group("/v3")
	{
		v3.GET("/ledgers", s.handler.QueryV3Ledgers)
		v3.GET("/ledgers/:height", s.handler.QueryV3Ledger)

		v3.GET("/transactions", s.handler.QueryV3Txs)
		v3.GET("/transactions/:txhash", s.handler.QueryV3SingleTx)
		v3.GET("/ledgers/:height/transactions", s.handler.QueryV3LedgerTxs)
		v3.GET("/accounts/:address/transactions", s.handler.QueryV3AccTxs)

		v3.GET("/payments", s.handler.QueryV3Payments)
		v3.GET("/ledgers/:height/payments", s.handler.QueryV3LedgerPayments)
		v3.GET("/accounts/:address/payments", s.handler.QueryV3AccPayments)
		v3.GET("/transactions/:txhash/payments", s.handler.QueryV3TxPayments)

		v3.GET("/config/tokens", s.handler.V3QueryConfigTokens)
		v3.GET("/config/nodes", s.handler.V3QueryConfigNodes)
		v3.GET("/ext/price/:symbol", cache.CachePageAtomic(store, time.Minute, s.handler.V3QueryPrice))
		//v3.GET("/ext/price/:symbol", s.handler.V3QueryPrice)
	}

	// reverse proxy
	s.proxy.SetPrefixPath("/v2/proxy")
	router.GET("/v2/proxy/*proxypath", gin.WrapH(s.proxy.Proxy()))

	if len(os.Args) > 1 && os.Args[1] == "version" {
		return
	}

	router.Run(s.cfg.Bind)
}
