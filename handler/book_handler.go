package handler

import (
	"net/http"

	"github.com/gadhittana-01/book-go/constant"
	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana-01/book-go/service"
	"github.com/gadhittana-01/book-go/utils"
	"github.com/go-chi/chi"
)

type BookHandler interface {
	SetupBookRoutes(route *chi.Mux)
}

type BookHandlerImpl struct {
	bookSvc        service.BookSvc
	authMiddleware utils.AuthMiddleware
}

func NewBookHandler(
	bookSvc service.BookSvc,
	authMiddleware utils.AuthMiddleware,
) BookHandler {
	return &BookHandlerImpl{
		bookSvc:        bookSvc,
		authMiddleware: authMiddleware,
	}
}

func (h *BookHandlerImpl) SetupBookRoutes(route *chi.Mux) {
	setupBookV1Routes(route, h)
}

func (h *BookHandlerImpl) CreateBook(w http.ResponseWriter, r *http.Request) {
	input := utils.ValidateBodyPayload(r.Body, &dto.CreateBookReq{})

	resp := h.bookSvc.CreateBook(r.Context(), input)

	utils.GenerateSuccessResp(w, resp, http.StatusCreated)
}

func (h *BookHandlerImpl) GetBook(w http.ResponseWriter, r *http.Request) {
	page := utils.ValidateQueryParamInt(r, "page", 1)
	limit := utils.ValidateQueryParamInt(r, "limit", constant.DefaultLimit)

	resp := h.bookSvc.GetBook(r.Context(), dto.GetBookReq{
		Page:  int32(page),
		Limit: int32(limit),
	})

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func setupBookV1Routes(route *chi.Mux, h *BookHandlerImpl) {
	route.Post("/v1/book", h.authMiddleware.CheckIsAuthenticated(h.CreateBook))
	route.Get("/v1/book", h.authMiddleware.CheckIsAuthenticated(h.GetBook))
}
