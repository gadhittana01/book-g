package service

import (
	"context"

	"github.com/gadhittana-01/book-go/constant"
	querier "github.com/gadhittana-01/book-go/db/repository"
	"github.com/gadhittana-01/book-go/dto"
	utilsConstant "github.com/gadhittana01/go-modules/constant"
	"github.com/gadhittana01/go-modules/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

const (
	FailedToFindBookByID             = "Failed to find book by ID"
	FailedToCheckBookExists          = "Failed to check book exists"
	FailedToCreateBook               = "Failed to create book"
	FailedToGetBook                  = "Failed to get book"
	FailedToGetBookPurchasedByUserID = "Failed to get book purchased by user id"
)

type (
	PaginationBookResp = dto.PaginationResp[dto.GetBookRes]
)

type BookSvc interface {
	CreateBook(ctx context.Context, input dto.CreateBookReq) dto.CreateBookRes
	GetBook(ctx context.Context, input dto.GetBookReq) PaginationBookResp
	GetBookPuchasedByUser(ctx context.Context) []dto.GetBookPuchasedByUserRes
}

type BookSvcImpl struct {
	repo     querier.Repository
	config   *utils.BaseConfig
	cacheSvc utils.CacheSvc
}

func NewBookSvc(
	repo querier.Repository,
	config *utils.BaseConfig,
	cacheSvc utils.CacheSvc,
) BookSvc {

	return &BookSvcImpl{
		repo:     repo,
		config:   config,
		cacheSvc: cacheSvc,
	}
}

func (s *BookSvcImpl) CreateBook(ctx context.Context, input dto.CreateBookReq) dto.CreateBookRes {
	var resp dto.CreateBookRes
	var book querier.Book
	var err error

	err = utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		book, err = repoTx.CreateBook(ctx, querier.CreateBookParams{
			Title:       input.Title,
			Description: input.Description,
			Author:      input.Author,
			Price:       input.Price,
		})
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToCreateBook, 422)
		}

		return nil
	})
	utils.PanicIfError(err)
	s.cacheSvc.ClearCaches([]string{constant.BookCacheKey}, "")

	resp = dto.CreateBookRes{
		ID:          book.ID.String(),
		Title:       book.Title,
		Description: book.Description,
		Author:      book.Author,
		Price:       book.Price,
	}

	return resp
}

func (s *BookSvcImpl) GetBook(ctx context.Context, input dto.GetBookReq) dto.PaginationResp[dto.GetBookRes] {
	resp, err := utils.GetOrSetData(s.cacheSvc, utils.BuildCacheKey(constant.BookCacheKey,
		"", "GetBook", input), func() (dto.PaginationResp[dto.GetBookRes], error) {
		ewg := errgroup.Group{}
		var err1 error
		var err2 error
		var books []querier.Book
		var count int64

		ewg.Go(func() error {
			books, err1 = s.repo.FindBook(ctx, querier.FindBookParams{
				Limit:  input.Limit,
				Offset: (input.Page - 1) * input.Limit,
			})
			return err1
		})

		ewg.Go(func() error {
			count, err2 = s.repo.GetBookCount(ctx)
			return err2
		})

		if err := ewg.Wait(); err != nil {
			return dto.PaginationResp[dto.GetBookRes]{}, utils.CustomErrorWithTrace(err,
				FailedToGetBook, 400)
		}

		return dto.ToPaginationResp(lo.Map(books, func(item querier.Book, index int) dto.GetBookRes {
			return dto.GetBookRes{
				ID:          item.ID.String(),
				Title:       item.Title,
				Description: item.Description,
				Author:      item.Author,
				Price:       item.Price,
			}
		}), int(input.Page), int(input.Limit), int(count)), nil
	})
	utils.PanicIfError(err)

	return resp
}

func (s *BookSvcImpl) GetBookPuchasedByUser(ctx context.Context) []dto.GetBookPuchasedByUserRes {
	authPayload := utils.GetRequestCtx(ctx, utilsConstant.UserSession)

	userID, err := uuid.Parse(authPayload.UserID)
	utils.PanicIfAppError(err, FailedToParseStringToUUID, 400)

	res, err := s.repo.GetBookPurchasedByUserID(ctx, userID)
	utils.PanicIfAppError(err, FailedToGetBookPurchasedByUserID, 400)

	return lo.Map(res, func(item querier.GetBookPurchasedByUserIDRow, index int) dto.GetBookPuchasedByUserRes {
		return dto.GetBookPuchasedByUserRes{
			ID:          item.BookID.String(),
			Title:       item.Title,
			Description: item.Description,
		}
	})
}
