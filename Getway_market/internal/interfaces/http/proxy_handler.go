package http

import (
	"context"
	"net/http"
	"time"

	"Getway_market/internal/application"
	"Getway_market/internal/infrastructure/proxy"
)

type ProxyHandler struct {
	router *application.RouterService
}

func NewProxyHandler(router *application.RouterService) *ProxyHandler {
	return &ProxyHandler{router: router}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// bearer: из Authorization, иначе попробуем cookie Access_Token
	authz := r.Header.Get("Authorization")
	if authz == "" {
		if c, err := r.Cookie("Access_Token"); err == nil && c.Value != "" {
			authz = "Bearer " + c.Value
		}
	}

	res, err := h.router.Resolve(
		r.Context(),
		r.Method,
		r.URL.Path,
		r.URL.Query(),
		authz,
	)
	if err != nil {
		switch err {
		case application.ErrUnauthorized:
			http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"invalid token"}}`, http.StatusUnauthorized)
		case application.ErrNoRoute:
			http.Error(w, `{"error":{"code":"NOT_FOUND","message":"no route"}}`, http.StatusNotFound)
		case application.ErrNoRule:
			http.Error(w, `{"error":{"code":"NOT_FOUND","message":"no rule"}}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error":{"code":"BAD_GATEWAY","message":"resolve error"}}`, http.StatusBadGateway)
		}
		return
	}

	// таймаут на апстрим
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(res.TimeoutMs)*time.Millisecond)
	defer cancel()
	r = r.WithContext(ctx)

	// переписываем путь + query
	r.URL.Path = res.Path
	r.URL.RawQuery = res.Query.Encode()

	// директор внутри proxy установит scheme/host
	p := proxy.NewReverseProxy(res.TargetBase)
	p.ServeHTTP(w, r)
}
