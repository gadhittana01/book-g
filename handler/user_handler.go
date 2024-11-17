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

// SignUp godoc
// @Id signUp
// @Summary      Sign Up
// @Description  Sign Up
// @Tags         auth
// @Accept 		 json
// @Param		 requestBody		body		dto.SignUpReq	true	"Sign Up Request"
// @Produce      json
// @Success      200  {object}  dto.SuccessResp200{data=dto.SignUpRes}
// @Failure      400  {object}  dto.FailedResp400
// @Failure      401  {object}  dto.FailedResp401
// @Failure      404  {object}  dto.FailedResp404
// @Failure      500  {object}  dto.FailedResp500
// @Router       /v1/sign-up [post]
func (h *UserHandlerImpl) SignUp(w http.ResponseWriter, r *http.Request) {
	input := utils.ValidateBodyPayload(r.Body, &dto.SignUpReq{})

	resp := h.userSvc.SignUp(r.Context(), input)

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

// SignIn godoc
// @Id signIn
// @Summary      Sign In
// @Description  Sign In
// @Tags         auth
// @Accept 		 json
// @Param		 requestBody		body		dto.SignInReq	true	"Sign In Request"
// @Produce      json
// @Success      200  {object}  dto.SuccessResp200{data=dto.SignInRes}
// @Failure      400  {object}  dto.FailedResp400
// @Failure      401  {object}  dto.FailedResp401
// @Failure      404  {object}  dto.FailedResp404
// @Failure      500  {object}  dto.FailedResp500
// @Router       /v1/sign-in [post]
func (h *UserHandlerImpl) SignIn(w http.ResponseWriter, r *http.Request) {
	input := utils.ValidateBodyPayload(r.Body, &dto.SignInReq{})

	resp := h.userSvc.SignIn(r.Context(), input)

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func (h *UserHandlerImpl) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.GenerateSuccessResp(w, "UP", http.StatusOK)
}

func setupUserV1Routes(route *chi.Mux, h *UserHandlerImpl) {
	route.Get("/v1/health-check", h.HealthCheck)
	route.Post("/v1/sign-up", h.SignUp)
	route.Post("/v1/sign-in", h.SignIn)
}
