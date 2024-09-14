package handler

import (
	"net/http"

	"github.com/gadhittana-01/book-go/constant"
	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana-01/book-go/service"
	"github.com/gadhittana-01/book-go/utils"
	"github.com/go-chi/chi"
)

const (
	InvalidOrderDetail = "Order detail cannot be empty"
)

type OrderHandler interface {
	SetupOrderRoutes(route *chi.Mux)
}

type OrderHandlerImpl struct {
	orderSvc       service.OrderSvc
	authMiddleware utils.AuthMiddleware
}

func NewOrderHandler(
	orderSvc service.OrderSvc,
	authMiddleware utils.AuthMiddleware,
) OrderHandler {
	return &OrderHandlerImpl{
		orderSvc:       orderSvc,
		authMiddleware: authMiddleware,
	}
}

func (h *OrderHandlerImpl) SetupOrderRoutes(route *chi.Mux) {
	setupOrderV1Routes(route, h)
}

func (h *OrderHandlerImpl) CreateOrder(w http.ResponseWriter, r *http.Request) {
	input := utils.ValidateBodyPayload(r.Body, &dto.CreateOrderReq{})

	if len(input.OrderDetail) == 0 {
		utils.PanicAppError(InvalidOrderDetail, 400)
	}

	resp := h.orderSvc.CreateOrder(r.Context(), input)

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func (h *OrderHandlerImpl) GetOrder(w http.ResponseWriter, r *http.Request) {
	page := utils.ValidateQueryParamInt(r, "page", 1)
	limit := utils.ValidateQueryParamInt(r, "limit", constant.DefaultLimit)

	resp := h.orderSvc.GetOrder(r.Context(), dto.GetOrderReq{
		Page:  int32(page),
		Limit: int32(limit),
	})

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func (h *OrderHandlerImpl) GetOrderDetail(w http.ResponseWriter, r *http.Request) {
	orderID := utils.ValidateURLParamUUID(r, "orderId")

	resp := h.orderSvc.GetOrderDetail(r.Context(), dto.GetOrderDetailReq{
		OrderID: orderID,
	})

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func setupOrderV1Routes(route *chi.Mux, h *OrderHandlerImpl) {
	route.Post("/v1/order", h.authMiddleware.CheckIsAuthenticated(h.CreateOrder))
	route.Get("/v1/order", h.authMiddleware.CheckIsAuthenticated(h.GetOrder))
	route.Get("/v1/order/{orderId}", h.authMiddleware.CheckIsAuthenticated(h.GetOrderDetail))
}
