package app

import (
	"testing"
	"time"

	"github.com/gadhittana-01/book-go/handler"
	mocksvc "github.com/gadhittana-01/book-go/service/mock"
	"github.com/gadhittana01/go-modules/utils"
	mockutl "github.com/gadhittana01/go-modules/utils/mock"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func initApp(ctrl *gomock.Controller) App {
	r := chi.NewRouter()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)

	mockToken := mockutl.NewMockTokenClient(ctrl)
	authMiddleware := utils.NewAuthMiddleware(config, mockToken)
	userSvc := mocksvc.NewMockUserSvc(ctrl)
	orderSvc := mocksvc.NewMockOrderSvc(ctrl)
	bookSvc := mocksvc.NewMockBookSvc(ctrl)
	userHandler := handler.NewUserHandler(userSvc)
	orderHandler := handler.NewOrderHandler(orderSvc, authMiddleware)
	bookHandler := handler.NewBookHandler(bookSvc, authMiddleware)

	return NewApp(r, config, userHandler, orderHandler, bookHandler)
}

func TestNewApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := initApp(ctrl)
	assert.NotNil(t, app)
	assert.IsType(t, &AppImpl{}, app)
}

func TestStartApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := initApp(ctrl)

	// START SERVER ON BACKGROUND
	go app.Start()
	time.Sleep(500 * time.Millisecond)
}
