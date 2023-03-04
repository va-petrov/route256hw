package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/lomsclient"
	"route256/checkout/internal/clients/productsclient"
	"route256/checkout/internal/config"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/checkout/internal/service"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lomsClient := lomsclient.New(config.ConfigData.Services.Loms)
	productsClient := productsclient.New(config.ConfigData.Services.ProductService.Url,
		config.ConfigData.Services.ProductService.Token)
	defer productsClient.Close()

	checkoutService := service.New(lomsClient, productsClient)

	purchaseHandler := purchase.New(checkoutService)
	http.Handle("/purchase", srvwrapper.New(purchaseHandler.Handle))
	addToCartHandler := addtocart.New(checkoutService)
	http.Handle("/addToCart", srvwrapper.New(addToCartHandler.Handle))
	deleteFromCartHandler := deletefromcart.New(checkoutService)
	http.Handle("/deleteFromCart", srvwrapper.New(deleteFromCartHandler.Handle))
	listHandler := listcart.New(checkoutService)
	http.Handle("/listCart", srvwrapper.New(listHandler.Handle))

	log.Println("listening http at", port)
	err = http.ListenAndServe(port, nil)
	log.Fatal("cannot listen http", err)
}
