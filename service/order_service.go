package service

import (
	"context"
	"time"

	"github.com/gadhittana-01/book-go/constant"
	querier "github.com/gadhittana-01/book-go/db/repository"
	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana-01/book-go/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

const (
	FailedToParseStringToUUID = "Failed to parse string to UUID"
	FailedToCreateOrder       = "Failed to create order"
	FailedToCreateOrderDetail = "Failed to create order detail"
	FailedToCheckBookExists   = "Failed to check book exists"
	FailedToGetOrder          = "Failed to get order"
	FailedToCheckOrderExists  = "Failed to check order exists"
	FailedToFindOrderByID     = "Failed to find order by ID"
	BookNotExists             = "Book doesn't exists"
	OrderNotExists            = "Order doesn't exists"
	InvalidBookID             = "BookID must UUID and cannot be empty"
	InvalidQuantity           = "Quantity must greater than zero"
)

type OrderSvc interface {
	CreateOrder(ctx context.Context, input dto.CreateOrderReq) dto.CreateOrderRes
	GetOrder(ctx context.Context, input dto.GetOrderReq) dto.PaginationResp[dto.GetOrderRes]
	GetOrderDetail(ctx context.Context, input dto.GetOrderDetailReq) dto.GetOrderDetailRes
}

type OrderSvcImpl struct {
	repo     querier.Repository
	config   *utils.BaseConfig
	token    utils.TokenClient
	cacheSvc utils.CacheSvc
}

func NewOrderSvc(
	repo querier.Repository,
	config *utils.BaseConfig,
	token utils.TokenClient,
	cacheSvc utils.CacheSvc,
) OrderSvc {
	return &OrderSvcImpl{
		repo:     repo,
		config:   config,
		token:    token,
		cacheSvc: cacheSvc,
	}
}

func (s *OrderSvcImpl) CreateOrder(ctx context.Context, input dto.CreateOrderReq) dto.CreateOrderRes {
	var resp dto.CreateOrderRes
	var order querier.Order
	now := time.Now()
	authPayload := utils.GetRequestCtx(ctx, constant.UserSession)

	userID, err := uuid.Parse(authPayload.UserID)
	utils.PanicIfAppError(err, FailedToParseStringToUUID, 400)

	err = utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		order, err = repoTx.CreateOrder(ctx, querier.CreateOrderParams{
			UserID: userID,
			Date:   now,
		})
		utils.PanicIfAppError(err, FailedToCreateOrder, 422)

		for _, item := range input.OrderDetail {
			if item.BookID == "" {
				utils.PanicAppError(InvalidBookID, 400)
			}

			if item.Quantity <= 0 {
				utils.PanicAppError(InvalidQuantity, 400)
			}

			bookID, err := uuid.Parse(item.BookID)
			utils.PanicIfAppError(err, FailedToParseStringToUUID, 400)

			isExists, err := repoTx.CheckBookExists(ctx, bookID)
			utils.PanicIfAppError(err, FailedToCheckBookExists, 400)

			if !isExists {
				utils.PanicAppError(BookNotExists, 400)
			}

			_, err = repoTx.CreateOrderDetail(ctx, querier.CreateOrderDetailParams{
				OrderID:  order.ID,
				BookID:   bookID,
				Quantity: int32(item.Quantity),
			})
			utils.PanicIfAppError(err, FailedToCreateOrderDetail, 422)
		}

		return nil
	})
	utils.PanicIfError(err)

	resp = dto.CreateOrderRes{
		OrderId: order.ID.String(),
		Date:    order.Date.Format(constant.TimeFormat),
	}

	return resp
}

func (s *OrderSvcImpl) GetOrder(ctx context.Context, input dto.GetOrderReq) dto.PaginationResp[dto.GetOrderRes] {
	authPayload := utils.GetRequestCtx(ctx, constant.UserSession)

	userID, err := uuid.Parse(authPayload.UserID)
	utils.PanicIfAppError(err, FailedToParseStringToUUID, 400)

	resp, err := utils.GetOrSetData(s.cacheSvc, utils.BuildCacheKey(constant.OrderCacheKey,
		authPayload.UserID, "GetOrder", input), func() (dto.PaginationResp[dto.GetOrderRes], error) {
		ewg := errgroup.Group{}
		var err1 error
		var err2 error
		var orders []querier.Order
		var count int64

		ewg.Go(func() error {
			orders, err1 = s.repo.FindOrderByUserID(ctx, querier.FindOrderByUserIDParams{
				UserID: userID,
				Limit:  input.Limit,
				Offset: (input.Page - 1) * input.Limit,
			})
			return err1
		})

		ewg.Go(func() error {
			count, err2 = s.repo.GetOrderCountByUserId(ctx, userID)
			return err2
		})

		if err := ewg.Wait(); err != nil {
			return dto.PaginationResp[dto.GetOrderRes]{}, utils.CustomErrorWithTrace(err,
				FailedToGetOrder, 400)
		}

		return dto.ToPaginationResp(lo.Map(orders, func(item querier.Order, index int) dto.GetOrderRes {
			return dto.GetOrderRes{
				OrderId: item.ID.String(),
				Date:    item.Date.Format(constant.TimeFormat),
			}
		}), int(input.Page), int(input.Limit), int(count)), err
	})
	utils.PanicIfError(err)

	return resp
}

func (s *OrderSvcImpl) GetOrderDetail(ctx context.Context, input dto.GetOrderDetailReq) dto.GetOrderDetailRes {
	authPayload := utils.GetRequestCtx(ctx, constant.UserSession)

	userID, err := uuid.Parse(authPayload.UserID)
	utils.PanicIfAppError(err, FailedToParseStringToUUID, 400)

	resp, err := utils.GetOrSetData(s.cacheSvc, utils.BuildCacheKey(constant.OrderCacheKey,
		authPayload.UserID, "GetOrderDetail", input), func() (dto.GetOrderDetailRes, error) {

		isExists, err := s.repo.CheckOrderExists(ctx, input.OrderID)
		if err != nil {
			return dto.GetOrderDetailRes{}, utils.CustomErrorWithTrace(err, FailedToCheckOrderExists, 400)
		}

		if !isExists {
			return dto.GetOrderDetailRes{}, utils.CustomError(OrderNotExists, 400)
		}

		order, err := s.repo.FindOrderByID(ctx, querier.FindOrderByIDParams{
			UserID: userID,
			ID:     input.OrderID,
		})
		if err != nil {
			return dto.GetOrderDetailRes{}, utils.CustomErrorWithTrace(err, FailedToFindOrderByID, 400)
		}

		orderDetail, err := s.repo.FindOrderDetailByOrderID(ctx, querier.FindOrderDetailByOrderIDParams{
			UserID: order.UserID,
			ID:     order.ID,
		})
		if err != nil {
			return dto.GetOrderDetailRes{}, utils.CustomErrorWithTrace(err, FailedToGetOrder, 400)
		}

		return dto.GetOrderDetailRes{
			OrderId: order.ID.String(),
			Date:    order.Date.Format(constant.TimeFormat),
			OrderDetail: lo.Map(orderDetail, func(item querier.FindOrderDetailByOrderIDRow, index int) dto.OrderDetail {
				return dto.OrderDetail{
					OrderDetailID: item.ID.String(),
					BookID:        item.BookID.String(),
					Title:         item.Title,
					Description:   item.Description,
					Author:        item.Author,
					Quantity:      int(item.Quantity),
				}
			}),
		}, nil
	})
	utils.PanicIfError(err)

	return resp
}
