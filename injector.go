//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gadhittana-01/book-go/app"
	querier "github.com/gadhittana-01/book-go/db/repository"
	"github.com/gadhittana-01/book-go/handler"
	"github.com/gadhittana-01/book-go/service"
	"github.com/gadhittana-01/book-go/utils"
	"github.com/go-chi/chi"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var userHandlerSet = wire.NewSet(
	querier.NewRepository,
	utils.NewToken,
	handler.NewUserHandler,
	service.NewUserSvc,
)

var orderHandlerSet = wire.NewSet(
	handler.NewOrderHandler,
	service.NewOrderSvc,
)

var authMiddlewareSet = wire.NewSet(
	utils.NewAuthMiddleware,
)

var cacheSet = wire.NewSet(
	wire.Bind(new(utils.RedisClient), new(*redis.Client)),
	utils.NewRedisClient,
	utils.NewCacheSvc,
)

func InitializeApp(
	route *chi.Mux,
	DB utils.PGXPool,
	config *utils.BaseConfig,
) (app.App, error) {
	wire.Build(
		userHandlerSet,
		orderHandlerSet,
		cacheSet,
		authMiddlewareSet,
		app.NewApp,
	)

	return nil, nil
}
