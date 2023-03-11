package loms_v1

import (
	"route256/loms/internal/service"
	"route256/loms/pkg/loms_v1"
)

type Implementation struct {
	loms_v1.UnimplementedLOMSServiceServer

	lomsService *service.Service
}

func NewLOMSV1(lomsService *service.Service) *Implementation {
	return &Implementation{
		loms_v1.UnimplementedLOMSServiceServer{},
		lomsService,
	}
}
