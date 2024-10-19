package handler

import (
	"net/http"

	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana-01/book-go/service"
	"github.com/gadhittana01/go-modules/utils"
	"github.com/go-chi/chi"
)

type UserHandler interface {
	SetupUserRoutes(route *chi.Mux)
}

type UserHandlerImpl struct {
	userSvc service.UserSvc
}

func NewUserHandler(
	userSvc service.UserSvc,
) UserHandler {
	return &UserHandlerImpl{
		userSvc: userSvc,
	}
}

func (h *UserHandlerImpl) SetupUserRoutes(route *chi.Mux) {
	setupUserV1Routes(route, h)
}

func (h *UserHandlerImpl) SignUp(w http.ResponseWriter, r *http.Request) {
	input := utils.ValidateBodyPayload(r.Body, &dto.SignUpReq{})

	resp := h.userSvc.SignUp(r.Context(), input)

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func (h *UserHandlerImpl) SignIn(w http.ResponseWriter, r *http.Request) {
	input := utils.ValidateBodyPayload(r.Body, &dto.SignInReq{})

	resp := h.userSvc.SignIn(r.Context(), input)

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func setupUserV1Routes(route *chi.Mux, h *UserHandlerImpl) {
	route.Post("/v1/sign-up", h.SignUp)
	route.Post("/v1/sign-in", h.SignIn)
}
