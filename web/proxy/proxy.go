package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

/*
反向代理包装
1. limitPath
2. upstreams
*/

type proxy struct {
	wrap *httputil.ReverseProxy
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.wrap.ServeHTTP(w, r)
}

type ReverseProxy struct {
	prefixPath string
	limitPath  []string
	upstreams  []string
}

func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{
		limitPath: make([]string, 0, 0),
		upstreams: make([]string, 0, 0),
	}
}

func (rp *ReverseProxy) Proxy() http.Handler {

	director := func(r *http.Request) {
		kill := func() {
			r.URL.Host = ""
			r.URL.Path = ""
			r.URL.Scheme = ""
			r.Host = ""
			r.RequestURI = ""
		}

		parts := strings.Split(r.URL.Path, rp.prefixPath)
		if len(parts) < 2 {
			kill()
			return
		}

		path := strings.Join(parts[1:], "/")
		if rp.isLimitPath(path) {
			kill()
			return
		}

		r.URL.Host = fmt.Sprintf("%s", rp.server())
		r.URL.Path = path
		r.URL.Scheme = "http"
		r.Host = r.URL.Host
	}

	return &proxy{
		wrap: &httputil.ReverseProxy{Director: director},
	}
}

func (rp *ReverseProxy) SetPrefixPath(path string) {
	rp.prefixPath = path
}

func (rp *ReverseProxy) AddToSetUpstream(servers ...string) {
	for _, server := range servers {
		isHit := false
		for _, v := range rp.upstreams {
			if server == v {
				isHit = true
				break
			}
		}
		if !isHit {
			rp.upstreams = append(rp.upstreams, server)
		}
	}
}

func (rp *ReverseProxy) AddToSetLimitPath(paths ...string) {
	for _, path := range paths {
		isHit := false
		for _, v := range rp.limitPath {
			if path == v {
				isHit = true
				break
			}
		}
		if !isHit {
			rp.limitPath = append(rp.limitPath, path)
		}
	}
}

func (rp *ReverseProxy) isLimitPath(path string) bool {
	for _, p := range rp.limitPath {
		//if strings.HasPrefix(path, p) {
		//	return true
		//}
		if p == path {
			return true
		}
	}
	return false
}

func (rp *ReverseProxy) server() string {
	if len(rp.upstreams) == 0 {
		return ""
	}

	// 目前就取第一个就行，当前应用程序只会取本地信息
	// 后期有需要再加上负载策略
	return rp.upstreams[0]
}
