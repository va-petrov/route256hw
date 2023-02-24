package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lomsClient := loms.New(config.ConfigData.Services.Loms)

	businessLogic := domain.New(lomsClient)

	addToCartHandler := addtocart.New(businessLogic)

	http.Handle("/addToCart", srvwrapper.New(addToCartHandler.Handle))

	log.Println("listening http at", port)
	err = http.ListenAndServe(port, nil)
	log.Fatal("cannot listen http", err)
}
