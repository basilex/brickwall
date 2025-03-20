package api

import (
	"context"

	_ "brickwall/docs"

	"brickwall/cmd/api/controller"
	"brickwall/internal/provider"

	swaggerDoc "github.com/swaggo/files"
	swaggerGin "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(ctx context.Context, router provider.IRouterProvider) {
	api := router.Engine().Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Swagger docs
			v1.GET("/docs/*any", swaggerGin.WrapHandler(swaggerDoc.Handler))

			// TODO: Auth middleware !!!

			// API controllers
			controller.NewAuxController(ctx, v1).Register()
			controller.NewUserController(ctx, v1).Register()
			controller.NewRoleController(ctx, v1).Register()
			controller.NewCountryController(ctx, v1).Register()
			controller.NewCurrencyController(ctx, v1).Register()
		}
	}
}
