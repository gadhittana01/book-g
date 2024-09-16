package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/gadhittana-01/book-go/constant"
	querier "github.com/gadhittana-01/book-go/db/repository"
	mockrepo "github.com/gadhittana-01/book-go/db/repository/mock"
	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana-01/book-go/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func initBookSvc(
	t *testing.T,
	ctrl *gomock.Controller,
	config *utils.BaseConfig,
) (BookSvc, *mockrepo.MockRepository, utils.CacheSvc) {
	mockRepo := mockrepo.NewMockRepository(ctrl)
	cacheSvc := utils.InitCacheSvc(t, config)
	return NewBookSvc(mockRepo, config, cacheSvc), mockRepo, cacheSvc
}

func TestCreateBook(t *testing.T) {
	ctx := utils.SetRequestContext(userID.String())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	bookSvcMock, mockRepo, _ := initBookSvc(t, ctrl, config)

	bookID := uuid.New()
	title := "Hello"
	description := "World"
	author := "Giri Putra Adhittana"
	price := float64(10)
	now := time.Now()
	req := dto.CreateBookReq{
		Title:       title,
		Description: description,
		Author:      author,
		Price:       price,
	}

	t.Run("success create Book", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo)

		mockRepo.EXPECT().CreateBook(gomock.Any(), querier.CreateBookParams{
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
		}).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, nil).Times(1)

		resp := bookSvcMock.CreateBook(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.CreateBookRes{
			ID:          bookID.String(),
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
		}, resp)
	})

	t.Run("failed to create Book", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CreateBook(gomock.Any(), querier.CreateBookParams{
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
		}).Return(querier.Book{
			ID:          bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 422,
			Message:    fmt.Sprintf("invalid request|%s", FailedToCreateBook),
		}, func() {
			resp := bookSvcMock.CreateBook(ctx, req)
			assert.Empty(t, resp)
		})
	})

}

func TestGetBook(t *testing.T) {
	ctx := utils.SetRequestContext(userID.String())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	bookSvcMock, mockRepo, mockCache := initBookSvc(t, ctrl, config)

	bookID := uuid.New()
	title := "Hello"
	description := "World"
	author := "Giri Putra Adhittana"
	price := float64(10)
	now := time.Now()
	page := int32(1)
	limit := int32(10)
	totalCount := 20
	req := dto.GetBookReq{
		Page:  page,
		Limit: limit,
	}

	t.Run("success get book", func(t *testing.T) {
		mockRepo.EXPECT().FindBook(gomock.Any(), querier.FindBookParams{
			Limit:  limit,
			Offset: (page - 1) * limit,
		}).Return([]querier.Book{
			{
				ID:          bookID,
				Title:       title,
				Description: description,
				Author:      author,
				Price:       price,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetBookCount(gomock.Any()).Return(int64(totalCount), nil).Times(1)

		resp := bookSvcMock.GetBook(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.PaginationResp[dto.GetBookRes]{
			Total:      totalCount,
			IsLoadMore: true,
			Next: dto.Next{
				Page: 2,
			},
			Prev: dto.Prev{
				Page: -1,
			},
			Data: []dto.GetBookRes{
				{
					ID:          bookID.String(),
					Title:       title,
					Description: description,
					Author:      author,
					Price:       price,
				},
			},
		}, resp)
	})

	t.Run("success get book from cache", func(t *testing.T) {
		resp := bookSvcMock.GetBook(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.PaginationResp[dto.GetBookRes]{
			Total:      totalCount,
			IsLoadMore: true,
			Next: dto.Next{
				Page: 2,
			},
			Prev: dto.Prev{
				Page: -1,
			},
			Data: []dto.GetBookRes{
				{
					ID:          bookID.String(),
					Title:       title,
					Description: description,
					Author:      author,
					Price:       price,
				},
			},
		}, resp)
	})

	t.Run("failed get book count", func(t *testing.T) {
		mockCache.DelByPrefix(ctx, constant.BookCacheKey)

		mockRepo.EXPECT().FindBook(gomock.Any(), querier.FindBookParams{
			Limit:  limit,
			Offset: (page - 1) * limit,
		}).Return([]querier.Book{
			{
				ID:          bookID,
				Title:       title,
				Description: description,
				Author:      author,
				Price:       price,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetBookCount(gomock.Any()).Return(int64(totalCount), errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToGetBook),
		}, func() {
			resp := bookSvcMock.GetBook(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed find book", func(t *testing.T) {
		mockCache.DelByPrefix(ctx, constant.BookCacheKey)

		mockRepo.EXPECT().FindBook(gomock.Any(), querier.FindBookParams{
			Limit:  limit,
			Offset: (page - 1) * limit,
		}).Return([]querier.Book{
			{
				ID:          bookID,
				Title:       title,
				Description: description,
				Author:      author,
				Price:       price,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}, errInvalidReq).Times(1)

		mockRepo.EXPECT().GetBookCount(gomock.Any()).Return(int64(totalCount), nil).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToGetBook),
		}, func() {
			resp := bookSvcMock.GetBook(ctx, req)
			assert.Empty(t, resp)
		})
	})
}
