package dto

import "github.com/google/uuid"

type SignUpReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type OrderDetailReq struct {
	BookID   string `json:"bookId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

type CreateOrderReq struct {
	OrderDetail []OrderDetailReq `json:"orderDetail" validate:"required"`
}

type GetOrderReq struct {
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
}

type GetOrderDetailReq struct {
	OrderID uuid.UUID `json:"orderId"`
}
