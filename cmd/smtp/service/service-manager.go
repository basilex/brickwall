package service

import (
	"context"
)

type IServiceManager interface {
	AuxService() IAuxService
}

type ServiceManager struct {
	ctx        context.Context
	auxService IAuxService
}

func NewServiceManager(ctx context.Context) IServiceManager {
	return &ServiceManager{
		ctx:        ctx,
		auxService: NewAuxService(ctx),
	}
}

func (rcv *ServiceManager) AuxService() IAuxService {
	return rcv.auxService
}
