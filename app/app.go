package app

import (
	"fmt"
	"net/http"

	"github.com/gadhittana-01/book-go/handler"
	"github.com/gadhittana-01/book-go/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type App interface {
	Start()
}

type AppImpl struct {
	route        *chi.Mux
	config       *utils.BaseConfig
	userHandler  handler.UserHandler
	orderHandler handler.OrderHandler
	bookHandler  handler.BookHandler
}

func NewApp(route *chi.Mux,
	config *utils.BaseConfig,
	userHandler handler.UserHandler,
	orderHandler handler.OrderHandler,
	bookHandler handler.BookHandler,
) App {
	return &AppImpl{
		route:        route,
		config:       config,
		userHandler:  userHandler,
		orderHandler: orderHandler,
		bookHandler:  bookHandler,
	}
}

func (s *AppImpl) Start() {
	s.route.Use(utils.Recovery)
	s.route.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	s.userHandler.SetupUserRoutes(s.route)
	s.orderHandler.SetupOrderRoutes(s.route)
	s.bookHandler.SetupBookRoutes(s.route)

	s.route.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.GenerateErrorResp[any](w, nil, 404)
	})

	utils.LogInfo(fmt.Sprintf("server started on port %d", s.config.ServerPort))
	port := fmt.Sprintf(":%d", s.config.ServerPort)
	err := http.ListenAndServe(port, s.route)
	if err != nil {
		panic(err)
	}
}
