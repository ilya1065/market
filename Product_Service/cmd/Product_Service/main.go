package main

import (
	"Product_Service/internal/config"
	"Product_Service/internal/inits"
	"Product_Service/internal/repository/product"
	"Product_Service/internal/service"
	https "Product_Service/internal/transport/http"
)

func main() {
	conf := config.NewConfig()
	db := inits.ConnectToDB(conf)
	repo := product.NewProdctRepo(db)
	servise := service.NewProductService(repo)
	hendler := https.NewHendler(*servise)
	router := https.NewRouter(hendler)
	router.Run(conf.HTTPPort)
}
