package main

import (
	"fmt"
	"log"
	"net/http"

	"Getway_market/internal/application"
	"Getway_market/internal/infrastructure/auth"
	"Getway_market/internal/infrastructure/config"
	"Getway_market/internal/infrastructure/registry"
	httpi "Getway_market/internal/interfaces/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("configs load error: ", err)
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	// services из yaml
	svcMap := map[string]string{}
	for name, s := range cfg.Routes.Services {
		svcMap[name] = s.BaseURL
	}
	reg := registry.NewMemoryRegistry(svcMap, cfg.Routes.Routes)
	authv := auth.NewJWTValidator(cfg.JWTSecret)
	router := application.NewRouterService(reg, authv)

	handler := httpi.NewProxyHandler(router)
	srv := httpi.NewServer(handler, cfg.CORSOrigins)

	log.Println("gateway listening on", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, srv.Handler); err != nil {
		log.Fatal(err)
	}
	fmt.Println("1234")
}
