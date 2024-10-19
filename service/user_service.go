package service

import (
	"context"

	querier "github.com/gadhittana-01/book-go/db/repository"
	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana01/go-modules/utils"
	"github.com/jackc/pgx/v5"
)

const (
	FailedToParseDate        = "Failed to parse date"
	FailedToCreateUser       = "Failed to create user"
	EmailAlreadyExist        = "Email already exist"
	FailedToFindUser         = "Failed to find user"
	FailedToCheckEmailExists = "Failed to check email exists"
	FailedToHashPassword     = "Failed to hash password"
	WrongCredentials         = "Wrong credentials"
	FailedToGenerateToken    = "Failed to generate token"
	UserNotFound             = "User not found"
)

type UserSvc interface {
	SignUp(ctx context.Context, input dto.SignUpReq) dto.SignUpRes
	SignIn(ctx context.Context, input dto.SignInReq) dto.SignInRes
}

type UserSvcImpl struct {
	repo   querier.Repository
	config *utils.BaseConfig
	token  utils.TokenClient
}

func NewUserSvc(
	repo querier.Repository,
	config *utils.BaseConfig,
	token utils.TokenClient,
) UserSvc {
	return &UserSvcImpl{
		repo:   repo,
		config: config,
		token:  token,
	}
}

func (s *UserSvcImpl) SignUp(ctx context.Context, input dto.SignUpReq) dto.SignUpRes {
	var resp dto.SignUpRes
	var user querier.User
	var token utils.GenerateTokenResp

	err := utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		isExists, err := repoTx.CheckEmailExists(ctx, input.Email)
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToCheckEmailExists, 400)
		}

		if isExists {
			return utils.CustomError(EmailAlreadyExist, 400)
		}

		pwd, err := utils.HashPassword(input.Password)
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToHashPassword, 400)
		}

		user, err = repoTx.CreateUser(ctx, querier.CreateUserParams{
			Name:     input.Name,
			Email:    input.Email,
			Password: pwd,
		})
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToCreateUser, 422)
		}

		token, err = s.token.GenerateToken(utils.GenerateTokenReq{
			UserID: user.ID.String(),
		})
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToGenerateToken, 400)
		}

		return nil
	})
	utils.PanicIfError(err)

	resp = dto.SignUpRes{
		ID:       user.ID.String(),
		Name:     user.Name,
		Token:    token.Token,
		ExpToken: token.ExpToken,
	}

	return resp
}

func (s *UserSvcImpl) SignIn(ctx context.Context, input dto.SignInReq) dto.SignInRes {
	var resp dto.SignInRes
	var user querier.User
	var err error

	user, err = s.repo.FindUserByEmail(ctx, input.Email)
	if err != nil && err != pgx.ErrNoRows {
		utils.PanicIfAppError(err, FailedToFindUser, 400)
	}

	if err == pgx.ErrNoRows {
		utils.PanicIfAppError(err, UserNotFound, 400)
	}

	if !utils.IsCorrectPassword(input.Password, user.Password) {
		utils.PanicAppError(WrongCredentials, 400)
	}

	token, err := s.token.GenerateToken(utils.GenerateTokenReq{
		UserID: user.ID.String(),
	})
	utils.PanicIfAppError(err, FailedToGenerateToken, 400)

	resp = dto.SignInRes{
		ID:       user.ID.String(),
		Name:     user.Name,
		Token:    token.Token,
		ExpToken: token.ExpToken,
	}

	return resp
}
