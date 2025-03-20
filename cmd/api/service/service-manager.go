package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"brickwall/internal/common"
	"brickwall/internal/storage/dbs"
)

type IServiceManager interface {
	AuxService() IAuxService
	UserService() IUserService
	AuthService() IAuthService
	RoleService() IRoleService
	CountryService() ICountryService
	CurrencyService() ICurrencyService
}
type ServiceManager struct {
	ctx     context.Context
	pool    *pgxpool.Pool
	queries *dbs.Queries

	auxService      IAuxService
	userService     IUserService
	authService     IAuthService
	roleService     IRoleService
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
		userService:     NewUserService(ctx, queries),
		authService:     NewAuthService(ctx, queries),
		roleService:     NewRoleService(ctx, queries),
		countryService:  NewCountryService(ctx, queries),
		currencyService: NewCurrencyService(ctx, queries),
	}
}

func (rcv *ServiceManager) AuxService() IAuxService {
	return rcv.auxService
}

func (rcv *ServiceManager) UserService() IUserService {
	return rcv.userService
}

func (rcv *ServiceManager) AuthService() IAuthService {
	return rcv.authService
}

func (rcv *ServiceManager) RoleService() IRoleService {
	return rcv.roleService
}

func (rcv *ServiceManager) CountryService() ICountryService {
	return rcv.countryService
}

func (rcv *ServiceManager) CurrencyService() ICurrencyService {
	return rcv.currencyService
}
