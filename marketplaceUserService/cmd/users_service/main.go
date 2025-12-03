package main

import (
	"log"
	"marketplace/internal/auth"
	"marketplace/internal/config"
	"marketplace/internal/repository"
	"marketplace/internal/service"
	httptr "marketplace/internal/transport/http"
	"marketplace/migrations"
)

func init() {
	migrations.ConnectToDB()

}

func main() {
	migrations.Migration()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewUserRepo(migrations.DB)
	j := auth.NewJWT(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	svc := service.NewUserService(repo, j)
	h := httptr.NewHandler(svc)
	r := httptr.NewRouter(h)

	log.Println("Users service listening on", cfg.HTTPPort)
	if err = r.Run(cfg.HTTPPort); err != nil {
		log.Fatal(err)
	}
}
