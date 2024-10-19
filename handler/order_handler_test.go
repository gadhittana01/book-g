package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gadhittana-01/book-go/constant"
	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana-01/book-go/service"
	mocksvc "github.com/gadhittana-01/book-go/service/mock"
	"github.com/gadhittana01/go-modules/utils"
	mockutl "github.com/gadhittana01/go-modules/utils/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderMock := mocksvc.NewMockOrderSvc(ctrl)
	middlewareMock := mockutl.NewMockAuthMiddleware(ctrl)

	type args struct {
		service        service.OrderSvc
		authMiddleware utils.AuthMiddleware
	}

	tests := []struct {
		name string
		args args
		want *OrderHandlerImpl
	}{
		{
			args: args{
				service:        orderMock,
				authMiddleware: middlewareMock,
			},
			want: &OrderHandlerImpl{
				orderSvc:       orderMock,
				authMiddleware: middlewareMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderHandler(tt.args.service, tt.args.authMiddleware); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderID := uuid.New()
	now := time.Now()
	totalPrice := float64(100)
	status := "pending"
	bookID := uuid.New()
	quantity := 10

	sampleReq := httptest.NewRequest("POST", "http://localhost:8000/v1/order", strings.NewReader(fmt.Sprintf(`{
		"orderDetail" : [
			{
				"bookId" : "%s",
				"quantity" : %d
			}
		]
	}`, bookID.String(), quantity)))
	sampleResp := httptest.NewRecorder()

	invalidSampleReq := httptest.NewRequest("POST", "http://localhost:8000/v1/order", strings.NewReader(`{
		"orderDetail" : [
		]
	}`))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.OrderSvc
	}

	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}

	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
	}{
		{
			name: "success create order",
			fields: func() fields {
				orderMock := mocksvc.NewMockOrderSvc(ctrl)

				orderMock.EXPECT().CreateOrder(gomock.Any(), dto.CreateOrderReq{
					OrderDetail: []dto.OrderDetailReq{
						{
							BookID:   bookID.String(),
							Quantity: quantity,
						},
					},
				}).Return(dto.CreateOrderRes{
					OrderId:    orderID.String(),
					Date:       now.Format(constant.TimeFormat),
					TotalPrice: totalPrice,
					Status:     status,
				}).Times(1)

				return fields{
					service: orderMock,
				}
			},
			args: args{
				w:   sampleResp,
				req: sampleReq,
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			fields: func() fields {
				orderMock := mocksvc.NewMockOrderSvc(ctrl)

				orderMock.EXPECT().CreateOrder(gomock.Any(), dto.CreateOrderReq{
					OrderDetail: []dto.OrderDetailReq{
						{
							BookID:   bookID.String(),
							Quantity: quantity,
						},
					},
				}).Return(dto.CreateOrderRes{
					OrderId:    orderID.String(),
					Date:       now.Format(constant.TimeFormat),
					TotalPrice: totalPrice,
					Status:     status,
				}).Times(0)

				return fields{
					service: orderMock,
				}
			},
			args: args{
				w:   invalidSampleResp,
				req: invalidSampleReq,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.fields()
			i := OrderHandlerImpl{
				orderSvc: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.CreateOrder(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.CreateOrder(tt.args.w, tt.args.req)
				})
			}

		})
	}
}

func TestGetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderID := uuid.New()
	now := time.Now()
	totalPrice := float64(100)
	status := "pending"
	page := 1
	limit := 10
	totalCount := 20

	sampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/v1/order?page=%d&limit=%d", page, limit), strings.NewReader(``))
	sampleResp := httptest.NewRecorder()

	invalidSampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/v1/order?page=%d&limit=test", page), strings.NewReader(``))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.OrderSvc
	}

	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}

	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
	}{
		{
			name: "success get order",
			fields: func() fields {
				orderMock := mocksvc.NewMockOrderSvc(ctrl)

				orderMock.EXPECT().GetOrder(gomock.Any(), dto.GetOrderReq{
					Page:  int32(page),
					Limit: int32(limit),
				}).Return(dto.PaginationResp[dto.GetOrderRes]{
					Total:      totalCount,
					IsLoadMore: true,
					Data: []dto.GetOrderRes{
						{
							OrderId:    orderID.String(),
							Date:       now.Format(constant.TimeFormat),
							TotalPrice: totalPrice,
							Status:     status,
						},
					},
				}).Times(1)

				return fields{
					service: orderMock,
				}
			},
			args: args{
				w:   sampleResp,
				req: sampleReq,
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			fields: func() fields {
				orderMock := mocksvc.NewMockOrderSvc(ctrl)

				orderMock.EXPECT().GetOrder(gomock.Any(), dto.GetOrderReq{
					Page:  int32(page),
					Limit: int32(limit),
				}).Return(dto.PaginationResp[dto.GetOrderRes]{
					Total:      totalCount,
					IsLoadMore: true,
					Data: []dto.GetOrderRes{
						{
							OrderId:    orderID.String(),
							Date:       now.Format(constant.TimeFormat),
							TotalPrice: totalPrice,
							Status:     status,
						},
					},
				}).Times(0)

				return fields{
					service: orderMock,
				}
			},
			args: args{
				w:   invalidSampleResp,
				req: invalidSampleReq,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.fields()
			i := OrderHandlerImpl{
				orderSvc: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.GetOrder(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.GetOrder(tt.args.w, tt.args.req)
				})
			}

		})
	}
}

func TestGetOrderDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderID := uuid.New()
	now := time.Now()
	totalPrice := float64(100)
	status := "pending"
	title := "Hello"
	description := "World"
	author := "Giri Putra Adhittana"
	quantity := 1
	orderDetailID := uuid.New()
	bookID := uuid.New()

	invalidSampleReq := httptest.NewRequest("GET", "http://localhost:8000/v1/order/123", strings.NewReader(``))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.OrderSvc
	}

	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}

	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
	}{
		{
			name: "success get order detail",
			fields: func() fields {
				orderMock := mocksvc.NewMockOrderSvc(ctrl)

				orderMock.EXPECT().GetOrderDetail(gomock.Any(), dto.GetOrderDetailReq{
					OrderID: orderID,
				}).Return(dto.GetOrderDetailRes{
					OrderId:    orderID.String(),
					Date:       now.Format(constant.TimeFormat),
					TotalPrice: totalPrice,
					Status:     status,
					OrderDetail: []dto.OrderDetail{
						{
							OrderDetailID: orderDetailID.String(),
							BookID:        bookID.String(),
							Title:         title,
							Description:   description,
							Author:        author,
							Quantity:      quantity,
						},
					},
				}).Times(0)

				return fields{
					service: orderMock,
				}
			},
			args: args{
				w:   invalidSampleResp,
				req: invalidSampleReq,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.fields()
			i := OrderHandlerImpl{
				orderSvc: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.GetOrderDetail(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.GetOrderDetail(tt.args.w, tt.args.req)
				})
			}

		})
	}
}
