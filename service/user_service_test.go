package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	querier "github.com/gadhittana-01/book-go/db/repository"
	mockrepo "github.com/gadhittana-01/book-go/db/repository/mock"
	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana-01/book-go/utils"
	mockutl "github.com/gadhittana-01/book-go/utils/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

var userID = uuid.New()
var errInvalidReq = errors.New("invalid request")

func initUserSvc(
	ctrl *gomock.Controller,
	config *utils.BaseConfig,
) (UserSvc, *mockrepo.MockRepository, *mockutl.MockTokenClient) {
	mockRepo := mockrepo.NewMockRepository(ctrl)
	mockToken := mockutl.NewMockTokenClient(ctrl)
	return NewUserSvc(mockRepo, config, mockToken), mockRepo, mockToken
}

func TestSignUp(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	userSvcMock, mockRepo, mockToken := initUserSvc(ctrl, config)

	name := "Giri Putra Adhittana"
	email := "test@gmail.com"
	password := "123"
	token := "dmytoken"
	expToken := int64(1000)
	req := dto.SignUpReq{
		Name:     name,
		Email:    email,
		Password: password,
	}

	t.Run("success sign up", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo)

		mockRepo.EXPECT().CheckEmailExists(gomock.Any(), email).Return(false, nil).Times(1)

		mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateUserParams{})).DoAndReturn(func(_ any, params querier.CreateUserParams) (querier.User, error) {
			assert.Equal(t, name, params.Name)
			assert.Equal(t, email, params.Email)

			return querier.User{
				ID:       userID,
				Name:     name,
				Email:    email,
				Password: password,
			}, nil
		}).Times(1)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, nil).Times(1)

		resp := userSvcMock.SignUp(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.SignUpRes{
			ID:       userID.String(),
			Name:     name,
			Token:    token,
			ExpToken: expToken,
		}, resp)
	})

	t.Run("failed to generate token", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CheckEmailExists(gomock.Any(), email).Return(false, nil).Times(1)

		mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateUserParams{})).DoAndReturn(func(_ any, params querier.CreateUserParams) (querier.User, error) {
			assert.Equal(t, name, params.Name)
			assert.Equal(t, email, params.Email)

			return querier.User{
				ID:       userID,
				Name:     name,
				Email:    email,
				Password: password,
			}, nil
		}).Times(1)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToGenerateToken),
		}, func() {
			resp := userSvcMock.SignUp(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed to create user", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CheckEmailExists(gomock.Any(), email).Return(false, nil).Times(1)

		mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateUserParams{})).DoAndReturn(func(_ any, params querier.CreateUserParams) (querier.User, error) {
			assert.Equal(t, name, params.Name)
			assert.Equal(t, email, params.Email)

			return querier.User{
				ID:       userID,
				Name:     name,
				Email:    email,
				Password: password,
			}, errInvalidReq
		}).Times(1)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 422,
			Message:    fmt.Sprintf("invalid request|%s", FailedToCreateUser),
		}, func() {
			resp := userSvcMock.SignUp(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("email already exists", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CheckEmailExists(gomock.Any(), email).Return(true, nil).Times(1)

		mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateUserParams{})).DoAndReturn(func(_ any, params querier.CreateUserParams) (querier.User, error) {
			assert.Equal(t, name, params.Name)
			assert.Equal(t, email, params.Email)

			return querier.User{
				ID:       userID,
				Name:     name,
				Email:    email,
				Password: password,
			}, errInvalidReq
		}).Times(0)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("%s|%s", EmailAlreadyExist, EmailAlreadyExist),
		}, func() {
			resp := userSvcMock.SignUp(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed to check email exists", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo, true)

		mockRepo.EXPECT().CheckEmailExists(gomock.Any(), email).Return(true, errInvalidReq).Times(1)

		mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateUserParams{})).DoAndReturn(func(_ any, params querier.CreateUserParams) (querier.User, error) {
			assert.Equal(t, name, params.Name)
			assert.Equal(t, email, params.Email)

			return querier.User{
				ID:       userID,
				Name:     name,
				Email:    email,
				Password: password,
			}, errInvalidReq
		}).Times(0)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToCheckEmailExists),
		}, func() {
			resp := userSvcMock.SignUp(ctx, req)
			assert.Empty(t, resp)
		})
	})
}

func TestSignIn(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	userSvcMock, mockRepo, mockToken := initUserSvc(ctrl, config)

	email := "test@gmail.com"
	name := "Giri Putra Adhittana"
	token := "dmytoken"
	expToken := int64(1000)
	password := "123"
	hashPassword, _ := utils.HashPassword(password)
	req := dto.SignInReq{
		Email:    email,
		Password: password,
	}

	t.Run("success sign in", func(t *testing.T) {
		mockRepo.EXPECT().FindUserByEmail(gomock.Any(), email).Return(querier.User{
			ID:       userID,
			Name:     name,
			Email:    email,
			Password: hashPassword,
		}, nil).Times(1)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, nil).Times(1)

		resp := userSvcMock.SignIn(ctx, req)

		assert.NotEmpty(t, resp)
		assert.Equal(t, dto.SignInRes{
			ID:       userID.String(),
			Name:     name,
			Token:    token,
			ExpToken: expToken,
		}, resp)
	})

	t.Run("failed to generate token", func(t *testing.T) {
		mockRepo.EXPECT().FindUserByEmail(gomock.Any(), email).Return(querier.User{
			ID:       userID,
			Name:     name,
			Email:    email,
			Password: hashPassword,
		}, nil).Times(1)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToGenerateToken),
		}, func() {
			resp := userSvcMock.SignIn(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("incorrect password", func(t *testing.T) {
		mockRepo.EXPECT().FindUserByEmail(gomock.Any(), email).Return(querier.User{
			ID:       userID,
			Name:     name,
			Email:    email,
			Password: password,
		}, nil).Times(1)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("%s|%s", WrongCredentials, WrongCredentials),
		}, func() {
			resp := userSvcMock.SignIn(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().FindUserByEmail(gomock.Any(), email).Return(querier.User{
			ID:       userID,
			Name:     name,
			Email:    email,
			Password: password,
		}, pgx.ErrNoRows).Times(1)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("%s|%s", pgx.ErrNoRows.Error(), UserNotFound),
		}, func() {
			resp := userSvcMock.SignIn(ctx, req)
			assert.Empty(t, resp)
		})
	})

	t.Run("failed to find user by email", func(t *testing.T) {
		mockRepo.EXPECT().FindUserByEmail(gomock.Any(), email).Return(querier.User{
			ID:       userID,
			Name:     name,
			Email:    email,
			Password: password,
		}, errInvalidReq).Times(1)

		mockToken.EXPECT().GenerateToken(utils.GenerateTokenReq{
			UserID: userID.String(),
		}).Return(utils.GenerateTokenResp{
			Token:    token,
			ExpToken: expToken,
		}, errInvalidReq).Times(0)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid request|%s", FailedToFindUser),
		}, func() {
			resp := userSvcMock.SignIn(ctx, req)
			assert.Empty(t, resp)
		})
	})

}
