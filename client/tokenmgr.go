package client

import (
	"fmt"
	"github.com/axengine/httpc"
	"github.com/toolglobal/api/libs/log"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type token struct {
	Name      string
	Symbol    string
	Contract  string
	Decimals  int
	CreatedAt string
}

type TokenMgr struct {
	baseURL string
	chainId string

	mutex        sync.RWMutex
	latestTokens map[string]token
}

func NewTokenMgr(source string, chainId string) *TokenMgr {
	return &TokenMgr{
		baseURL:      source,
		chainId:      chainId,
		latestTokens: make(map[string]token),
	}
}

func (m *TokenMgr) Start() {
	fn := func() {
		tm := time.NewTimer(time.Millisecond * 100)
		for {
			select {
			case <-tm.C:
				if err := m.sync(); err != nil {
					log.Logger.Error("sync tokens", zap.Error(err))
				}
				tm.Reset(time.Minute)
			}
		}
	}
	go fn()
}

func (m *TokenMgr) Sync() error {
	return m.sync()
}

func (m *TokenMgr) sync() error {
	tokens, err := m.fetchAllCoins()
	if err != nil {
		return err
	}

	mp := make(map[string]token)
	for _, v := range tokens {
		mp[v.Contract] = v
	}
	m.mutex.Lock()
	m.latestTokens = mp
	m.mutex.Unlock()
	return nil
}

func (m *TokenMgr) fetchAllCoins() ([]token, error) {
	coins := make([]token, 0)
	page := 1
	pageSize := 200

	type req struct {
		ChainId  string `json:"chainId"`
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Status   int    `json:"status"`
	}

	type resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			List []struct {
				Address       string `json:"address"`
				ChainId       string `json:"chainId"`
				Contract      string `json:"contract"`
				Point         int    `json:"point"`
				Created       int64  `json:"created"`
				Icon          string `json:"icon"`
				Name          string `json:"name"`
				SwapSymbol    string `json:"swapSymbol"`
				Symbol        string `json:"symbol"`
				TotalIssuance string `json:"totalIssuance"`
				WebSite       string `json:"webSite"`
				WhitePaper    string `json:"whitePaper"`
			}
		} `json:"data"`
	}

	fetch := func(r req) ([]token, error) {
		res := resp{}
		err := httpc.New(m.baseURL).Path("/v1/token/list").
			Query("page", strconv.FormatInt(int64(r.Page), 32)).
			Query("pageSize", strconv.FormatInt(int64(r.PageSize), 32)).
			Query("chainId", r.ChainId).
			Query("status", strconv.FormatInt(int64(r.Status), 32)).Get(&res, httpc.TypeApplicationJson)
		if err != nil {
			return nil, err
		}
		log.Logger.Info("sync", zap.Any("tokens", res))
		if res.Code != 0 {
			return nil, fmt.Errorf("code %d msg %s", res.Code, res.Msg)
		}

		coins := make([]token, 0)
		if res.Data.List != nil {
			for _, v := range res.Data.List {
				coins = append(coins, token{
					Name:      v.Name,
					Symbol:    v.Symbol,
					Contract:  v.Contract,
					Decimals:  v.Point,
					CreatedAt: "", // 暂无发行时间
				})
			}
		}
		return coins, nil
	}

	for {
		r := req{
			ChainId:  m.chainId,
			Page:     page,
			PageSize: pageSize,
			Status:   2,
		}
		cs, err := fetch(r)
		if err != nil {
			return nil, err
		}

		for _, c := range cs {
			coins = append(coins, c)
		}

		if len(cs) < pageSize {
			break
		}
		page++
	}
	return coins, nil
}

func (m *TokenMgr) Token(contract string) (token, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	tk, ok := m.latestTokens[contract]
	return tk, ok
}

func (m *TokenMgr) Tokens() []token {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	ls := make([]token, len(m.latestTokens))

	for _, v := range m.latestTokens {
		ls = append(ls, v)
	}
	return ls
}
