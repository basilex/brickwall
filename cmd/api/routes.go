package api

import (
	"context"

	"brickwall/cmd/api/controller"
	"brickwall/internal/provider"
	// "brickwall/bsp/cmd/api/domain/crud/currency"
)

func RegisterRoutes(ctx context.Context, router provider.IRouterProvider) {
	api := router.Engine().Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// TODO: Auth middleware !!!
			controller.NewAuxController(ctx, v1).Register()
			controller.NewUserController(ctx, v1).Register()
			controller.NewRoleController(ctx, v1).Register()
			controller.NewCountryController(ctx, v1).Register()
			controller.NewCurrencyController(ctx, v1).Register()
		}
	}
}
