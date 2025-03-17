package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"brickwall/internal/common"
	"brickwall/internal/storage/dbs"
)

type IServiceManager interface {
	AuxService() IAuxService
	AuthService() IAuthService
	CountryService() ICountryService
	CurrencyService() ICurrencyService
}
type ServiceManager struct {
	ctx     context.Context
	pool    *pgxpool.Pool
	queries *dbs.Queries

	auxService      IAuxService
	authService     IAuthService
	countryService  ICountryService
	currencyService ICurrencyService
}

func NewServiceManager(ctx context.Context) IServiceManager {
	pool := ctx.Value(common.KeyPgxProvider).(*pgxpool.Pool)
	queries := dbs.New(pool)

	return &ServiceManager{
		ctx:     ctx,
		pool:    pool,
		queries: queries,

		auxService:      NewAuxService(ctx, queries),
		authService:     NewAuthService(ctx, queries),
		countryService:  NewCountryService(ctx, queries),
		currencyService: NewCurrencyService(ctx, queries),
	}
}

func (rcv *ServiceManager) AuxService() IAuxService {
	return rcv.auxService
}

func (rcv *ServiceManager) AuthService() IAuthService {
	return rcv.authService
}

func (rcv *ServiceManager) CountryService() ICountryService {
	return rcv.countryService
}

func (rcv *ServiceManager) CurrencyService() ICurrencyService {
	return rcv.currencyService
}
