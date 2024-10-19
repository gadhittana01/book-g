package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/gadhittana-01/book-go/constant"
	querier "github.com/gadhittana-01/book-go/db/repository"
	mockrepo "github.com/gadhittana-01/book-go/db/repository/mock"
	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana01/go-modules/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func initOrderSvc(
	t *testing.T,
	ctrl *gomock.Controller,
	config *utils.BaseConfig,
) (OrderSvc, *mockrepo.MockRepository, utils.CacheSvc) {
	mockRepo := mockrepo.NewMockRepository(ctrl)
	cacheSvc := utils.InitCacheSvc(t, config)
	return NewOrderSvc(mockRepo, config, cacheSvc), mockRepo, cacheSvc
}

func TestCreateOrder(t *testing.T) {
	ctx := utils.SetRequestContext(userID.String())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	orderSvcMock, mockRepo, _ := initOrderSvc(t, ctrl, config)

	bookID := uuid.New()
	quantity := 10
	req := dto.CreateOrderReq{
		OrderDetail: []dto.OrderDetailReq{
			{
				BookID:   bookID.String(),
				Quantity: quantity,
			},
		},
	}
	invalidBookIDReq := dto.CreateOrderReq{
		OrderDetail: []dto.OrderDetailReq{
			{
				BookID:   "123",
				Quantity: quantity,
			},
		},
	}
	invalidQuantityReq := dto.CreateOrderReq{
		OrderDetail: []dto.OrderDetailReq{
			{
				BookID:   bookID.String(),
				Quantity: 0,
			},
		},
	}
	invalidEmptyBookIDReq := dto.CreateOrderReq{
		OrderDetail: []dto.OrderDetailReq{
			{
				BookID:   "",
				Quantity: quantity,
			},
		},
	}
	now := time.Now()
	totalPrice := float64(100)
	orderID := uuid.New()
	status := "pending"
	orderDetailID := uuid.New()
	title := "Hello"
	description := "World"
	author := "Giri Putra Adhittana"
	price := float64(10)

	t.Run("success create order", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(true, nil).Times(1)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, nil).Times(1)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, nil).Times(1)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, nil).Times(1)

		resp := orderSvcMock.CreateOrder(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.CreateOrderRes{
			OrderId:    orderID.String(),
			Date:       now.Format(constant.TimeFormat),
			TotalPrice: totalPrice,
			Status:     status,
		}, resp)
	})

	t.Run("failed update order by ID", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(true, nil).Times(1)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, nil).Times(1)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, nil).Times(1)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 422,
			Message:    fmt.Sprintf("invalid request|%s", FailedToUpdateOrder),
		}, func() {
			resp := orderSvcMock.CreateOrder(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed to find book by ID", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(true, nil).Times(1)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, nil).Times(1)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, errInvalidReq).Times(1)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToFindBookByID),
		}, func() {
			resp := orderSvcMock.CreateOrder(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed to create order detail", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(true, nil).Times(1)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, errInvalidReq).Times(1)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 422,
			Message:    fmt.Sprintf("invalid request|%s", FailedToCreateOrderDetail),
		}, func() {
			resp := orderSvcMock.CreateOrder(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("book not exists", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(false, nil).Times(1)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("%s|%s", BookNotExists, BookNotExists),
		}, func() {
			resp := orderSvcMock.CreateOrder(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed to check book exists", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(false, errInvalidReq).Times(1)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToCheckBookExists),
		}, func() {
			resp := orderSvcMock.CreateOrder(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed to parse bookID", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(false, errInvalidReq).Times(0)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid UUID length: 3|%s", FailedToParseStringToUUID),
		}, func() {
			resp := orderSvcMock.CreateOrder(ctx, invalidBookIDReq)
			assert.Empty(t, resp)
		})
	})

	t.Run("invalid quantity", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(false, errInvalidReq).Times(0)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("%s|%s", InvalidQuantity, InvalidQuantity),
		}, func() {
			resp := orderSvcMock.CreateOrder(ctx, invalidQuantityReq)
			assert.Empty(t, resp)
		})
	})

	t.Run("invalid bookID", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(false, errInvalidReq).Times(0)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("%s|%s", InvalidBookID, InvalidBookID),
		}, func() {
			resp := orderSvcMock.CreateOrder(ctx, invalidEmptyBookIDReq)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed to create order", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateOrderParams{})).DoAndReturn(func(_ any, params querier.CreateOrderParams) (querier.Order, error) {
			assert.Equal(t, userID, params.UserID)

			return querier.Order{
				ID:        orderID,
				UserID:    userID,
				Date:      now,
				Status:    status,
				CreatedAt: now,
				UpdatedAt: now,
			}, errInvalidReq
		}).Times(1)

		mockRepo.EXPECT().CheckBookExists(gomock.Any(), bookID).Return(false, errInvalidReq).Times(0)

		mockRepo.EXPECT().CreateOrderDetail(gomock.Any(), querier.CreateOrderDetailParams{
			OrderID:  orderID,
			BookID:   bookID,
			Quantity: int32(quantity),
		}).Return(querier.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			BookID:    bookID,
			Quantity:  int32(quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().FindBookByID(gomock.Any(), bookID).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().UpdateOrderByID(gomock.Any(), querier.UpdateOrderByIDParams{
			ID:         orderID,
			TotalPrice: totalPrice,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 422,
			Message:    fmt.Sprintf("invalid request|%s", FailedToCreateOrder),
		}, func() {
			resp := orderSvcMock.CreateOrder(ctx, invalidEmptyBookIDReq)
			assert.Empty(t, resp)
		})
	})
}

func TestGetOrder(t *testing.T) {
	ctx := utils.SetRequestContext(userID.String())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	orderSvcMock, mockRepo, mockCache := initOrderSvc(t, ctrl, config)

	page := int32(1)
	limit := int32(10)
	req := dto.GetOrderReq{
		Page:  page,
		Limit: limit,
	}
	now := time.Now()
	totalPrice := float64(100)
	status := "pending"
	orderID := uuid.New()
	totalCount := 20

	t.Run("success get order", func(t *testing.T) {
		mockRepo.EXPECT().FindOrderByUserID(gomock.Any(), querier.FindOrderByUserIDParams{
			UserID: userID,
			Limit:  limit,
			Offset: (page - 1) * limit,
		}).Return([]querier.Order{
			{
				ID:         orderID,
				UserID:     userID,
				Date:       now,
				TotalPrice: totalPrice,
				Status:     status,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetOrderCountByUserId(gomock.Any(), userID).Return(int64(totalCount), nil).Times(1)

		resp := orderSvcMock.GetOrder(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.PaginationResp[dto.GetOrderRes]{
			Total: totalCount,
			Next: dto.Next{
				Page: 2,
			},
			Prev: dto.Prev{
				Page: -1,
			},
			IsLoadMore: true,
			Data: []dto.GetOrderRes{
				{
					OrderId:    orderID.String(),
					Date:       now.Format(constant.TimeFormat),
					TotalPrice: totalPrice,
					Status:     status,
				},
			},
		}, resp)
	})

	t.Run("success get order from cache", func(t *testing.T) {
		resp := orderSvcMock.GetOrder(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.PaginationResp[dto.GetOrderRes]{
			Total: totalCount,
			Next: dto.Next{
				Page: 2,
			},
			Prev: dto.Prev{
				Page: -1,
			},
			IsLoadMore: true,
			Data: []dto.GetOrderRes{
				{
					OrderId:    orderID.String(),
					Date:       now.Format(constant.TimeFormat),
					TotalPrice: totalPrice,
					Status:     status,
				},
			},
		}, resp)
	})

	t.Run("failed get order count by user ID", func(t *testing.T) {
		mockCache.DelByPrefix(ctx, constant.OrderCacheKey)

		mockRepo.EXPECT().FindOrderByUserID(gomock.Any(), querier.FindOrderByUserIDParams{
			UserID: userID,
			Limit:  limit,
			Offset: (page - 1) * limit,
		}).Return([]querier.Order{
			{
				ID:         orderID,
				UserID:     userID,
				Date:       now,
				TotalPrice: totalPrice,
				Status:     status,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetOrderCountByUserId(gomock.Any(), userID).Return(int64(totalCount), errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToGetOrder),
		}, func() {
			resp := orderSvcMock.GetOrder(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed find order by user ID", func(t *testing.T) {
		mockCache.DelByPrefix(ctx, constant.OrderCacheKey)

		mockRepo.EXPECT().FindOrderByUserID(gomock.Any(), querier.FindOrderByUserIDParams{
			UserID: userID,
			Limit:  limit,
			Offset: (page - 1) * limit,
		}).Return([]querier.Order{
			{
				ID:         orderID,
				UserID:     userID,
				Date:       now,
				TotalPrice: totalPrice,
				Status:     status,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}, errInvalidReq).Times(1)

		mockRepo.EXPECT().GetOrderCountByUserId(gomock.Any(), userID).Return(int64(totalCount), nil).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToGetOrder),
		}, func() {
			resp := orderSvcMock.GetOrder(ctx, req)
			assert.Empty(t, resp)
		})
	})

}

func TestGetOrderDetail(t *testing.T) {
	ctx := utils.SetRequestContext(userID.String())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	orderSvcMock, mockRepo, mockCache := initOrderSvc(t, ctrl, config)

	orderID := uuid.New()
	orderDetailID := uuid.New()
	bookID := uuid.New()
	title := "Hello"
	description := "World"
	author := "Giri Putra Adhittana"
	quantity := 10
	price := float64(10)
	req := dto.GetOrderDetailReq{
		OrderID: orderID,
	}
	now := time.Now()
	totalPrice := float64(100)
	status := "pending"

	t.Run("success get order detail", func(t *testing.T) {
		mockRepo.EXPECT().CheckOrderExists(gomock.Any(), orderID).Return(true, nil).Times(1)

		mockRepo.EXPECT().FindOrderByID(gomock.Any(), querier.FindOrderByIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, nil).Times(1)

		mockRepo.EXPECT().FindOrderDetailByOrderID(gomock.Any(), querier.FindOrderDetailByOrderIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return([]querier.FindOrderDetailByOrderIDRow{
			{
				ID:          orderDetailID,
				Date:        now,
				BookID:      bookID,
				Title:       title,
				TotalPrice:  totalPrice,
				Status:      status,
				Description: description,
				Author:      author,
				Quantity:    int32(quantity),
				Price:       price,
			},
		}, nil).Times(1)

		resp := orderSvcMock.GetOrderDetail(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.GetOrderDetailRes{
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
		}, resp)
	})

	t.Run("success get order detail from cache", func(t *testing.T) {
		resp := orderSvcMock.GetOrderDetail(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.GetOrderDetailRes{
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
		}, resp)
	})

	t.Run("failed to get order detail by order ID", func(t *testing.T) {
		mockCache.DelByPrefix(ctx, constant.OrderCacheKey)

		mockRepo.EXPECT().CheckOrderExists(gomock.Any(), orderID).Return(true, nil).Times(1)

		mockRepo.EXPECT().FindOrderByID(gomock.Any(), querier.FindOrderByIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, nil).Times(1)

		mockRepo.EXPECT().FindOrderDetailByOrderID(gomock.Any(), querier.FindOrderDetailByOrderIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return([]querier.FindOrderDetailByOrderIDRow{
			{
				ID:          orderDetailID,
				Date:        now,
				BookID:      bookID,
				Title:       title,
				TotalPrice:  totalPrice,
				Status:      status,
				Description: description,
				Author:      author,
				Quantity:    int32(quantity),
				Price:       price,
			},
		}, errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToCreateOrderDetail),
		}, func() {
			resp := orderSvcMock.GetOrderDetail(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed to find order by ID", func(t *testing.T) {
		mockCache.DelByPrefix(ctx, constant.OrderCacheKey)

		mockRepo.EXPECT().CheckOrderExists(gomock.Any(), orderID).Return(true, nil).Times(1)

		mockRepo.EXPECT().FindOrderByID(gomock.Any(), querier.FindOrderByIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(1)

		mockRepo.EXPECT().FindOrderDetailByOrderID(gomock.Any(), querier.FindOrderDetailByOrderIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return([]querier.FindOrderDetailByOrderIDRow{
			{
				ID:          orderDetailID,
				Date:        now,
				BookID:      bookID,
				Title:       title,
				TotalPrice:  totalPrice,
				Status:      status,
				Description: description,
				Author:      author,
				Quantity:    int32(quantity),
				Price:       price,
			},
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToFindOrderByID),
		}, func() {
			resp := orderSvcMock.GetOrderDetail(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("order not exists", func(t *testing.T) {
		mockCache.DelByPrefix(ctx, constant.OrderCacheKey)

		mockRepo.EXPECT().CheckOrderExists(gomock.Any(), orderID).Return(false, nil).Times(1)

		mockRepo.EXPECT().FindOrderByID(gomock.Any(), querier.FindOrderByIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().FindOrderDetailByOrderID(gomock.Any(), querier.FindOrderDetailByOrderIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return([]querier.FindOrderDetailByOrderIDRow{
			{
				ID:          orderDetailID,
				Date:        now,
				BookID:      bookID,
				Title:       title,
				TotalPrice:  totalPrice,
				Status:      status,
				Description: description,
				Author:      author,
				Quantity:    int32(quantity),
				Price:       price,
			},
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("%s|%s", OrderNotExists, OrderNotExists),
		}, func() {
			resp := orderSvcMock.GetOrderDetail(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed check order exists", func(t *testing.T) {
		mockCache.DelByPrefix(ctx, constant.OrderCacheKey)

		mockRepo.EXPECT().CheckOrderExists(gomock.Any(), orderID).Return(false, errInvalidReq).Times(1)

		mockRepo.EXPECT().FindOrderByID(gomock.Any(), querier.FindOrderByIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return(querier.Order{
			ID:         orderID,
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, errInvalidReq).Times(0)

		mockRepo.EXPECT().FindOrderDetailByOrderID(gomock.Any(), querier.FindOrderDetailByOrderIDParams{
			UserID: userID,
			ID:     orderID,
		}).Return([]querier.FindOrderDetailByOrderIDRow{
			{
				ID:          orderDetailID,
				Date:        now,
				BookID:      bookID,
				Title:       title,
				TotalPrice:  totalPrice,
				Status:      status,
				Description: description,
				Author:      author,
				Quantity:    int32(quantity),
				Price:       price,
			},
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToCheckOrderExists),
		}, func() {
			resp := orderSvcMock.GetOrderDetail(ctx, req)
			assert.Empty(t, resp)
		})
	})

}
