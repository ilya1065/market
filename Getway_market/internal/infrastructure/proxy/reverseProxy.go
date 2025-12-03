package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var hopHeaders = map[string]struct{}{
	"Connection":        {},
	"Proxy-Connection":  {},
	"Keep-Alive":        {},
	"Transfer-Encoding": {},
	"TE":                {},
	"Trailer":           {},
	"Upgrade":           {},
}

func NewReverseProxy(targetBase *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = targetBase.Scheme
		req.URL.Host = targetBase.Host
		for h := range hopHeaders {
			req.Header.Del(h)
		}
		// X-Forwarded-*
		if ip, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			prior := req.Header.Get("X-Forwarded-For")
			if prior == "" {
				req.Header.Set("X-Forwarded-For", ip)
			} else {
				req.Header.Set("X-Forwarded-For", prior+", "+ip)
			}
		}
		if req.Header.Get("X-Forwarded-Proto") == "" {
			if req.TLS != nil {
				req.Header.Set("X-Forwarded-Proto", "https")
			} else {
				req.Header.Set("X-Forwarded-Proto", "http")
			}
		}
		if req.Header.Get("X-Forwarded-Host") == "" {
			req.Header.Set("X-Forwarded-Host", req.Host)
		}
	}

	return &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		ModifyResponse: func(resp *http.Response) error { return nil },
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, `{"error":{"code":"BAD_GATEWAY","message":"upstream error"}}`, http.StatusBadGateway)
		},
	}
}
