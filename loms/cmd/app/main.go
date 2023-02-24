package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/config"
	"route256/loms/internal/handlers/cancelorder"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/listorder"
	"route256/loms/internal/handlers/orderpayed"
	"route256/loms/internal/handlers/stocks"
	"route256/loms/internal/service"
)

const port = ":8081"

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lomsService := service.New()

	http.Handle("/createOrder", srvwrapper.New(createorder.New(lomsService).Handle))
	http.Handle("/listOrder", srvwrapper.New(listorder.New(lomsService).Handle))
	http.Handle("/orderPayed", srvwrapper.New(orderpayed.New(lomsService).Handle))
	http.Handle("/cancelOrder", srvwrapper.New(cancelorder.New(lomsService).Handle))
	http.Handle("/stocks", srvwrapper.New(stocks.New(lomsService).Handle))

	log.Println("listening http at", port)
	err = http.ListenAndServe(port, nil)
	log.Fatal("cannot listen http", err)
}
