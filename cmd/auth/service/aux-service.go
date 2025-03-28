package service

import (
	"context"

	"brickwall/internal/common"
	"brickwall/internal/provider"
)

type IAuxService interface {
	Index() *common.Message
	Health() *common.Message
	Metadata() *common.Metadata
	Environment() map[string]string
}
type AuxService struct {
	ctx context.Context
}

func NewAuxService(ctx context.Context) IAuxService {
	return &AuxService{ctx: ctx}
}

func (rcv *AuxService) Index() *common.Message {
	return &common.Message{
		Message: "Brickwall platform auth service",
	}
}

func (rcv *AuxService) Health() *common.Message {
	return &common.Message{
		Message: "Alive and kicking!",
	}
}

func (rcv *AuxService) Metadata() *common.Metadata {
	metadata := rcv.ctx.Value(common.KeyMetadata).(*common.Metadata)
	metadata.Service = "auth"
	return metadata
}

func (rcv *AuxService) Environment() map[string]string {
	env := rcv.ctx.Value(common.KeyEnvProvider).(provider.IEnvProvider)
	return env.Environment()
}
