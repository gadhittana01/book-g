package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gadhittana-01/book-go/dto"
	"github.com/gadhittana-01/book-go/service"
	mocksvc "github.com/gadhittana-01/book-go/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	userMock := mocksvc.NewMockUserSvc(ctrl)

	type args struct {
		service service.UserSvc
	}

	tests := []struct {
		name string
		args args
		want *UserHandlerImpl
	}{
		{
			args: args{
				service: userMock,
			},
			want: &UserHandlerImpl{
				userSvc: userMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserHandler(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := "Giri Putra Adhittana"
	email := "test@gmail.com"
	password := "123"
	userID := uuid.New()
	token := "123"
	expToken := int64(1000)

	sampleReq := httptest.NewRequest("POST", "http://localhost:8000/v1/sign-up", strings.NewReader(fmt.Sprintf(`{
		"name" : "%s",
		"email" : "%s",
		"password" : "%s"
	}`, name, email, password)))
	sampleResp := httptest.NewRecorder()

	invalidSampleReq := httptest.NewRequest("POST", "http://localhost:8000/v1/sign-up", strings.NewReader(fmt.Sprintf(`{
		"email" : "%s",
		"password" : "%s"
	}`, email, password)))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.UserSvc
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
			name: "success sign up",
			fields: func() fields {
				userMock := mocksvc.NewMockUserSvc(ctrl)

				userMock.EXPECT().SignUp(gomock.Any(), dto.SignUpReq{
					Name:     name,
					Email:    email,
					Password: password,
				}).Return(dto.SignUpRes{
					ID:       userID.String(),
					Name:     name,
					Token:    token,
					ExpToken: expToken,
				}).Times(1)

				return fields{
					service: userMock,
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
				userMock := mocksvc.NewMockUserSvc(ctrl)

				userMock.EXPECT().SignUp(gomock.Any(), dto.SignUpReq{
					Name:     name,
					Email:    email,
					Password: password,
				}).Return(dto.SignUpRes{
					ID:       userID.String(),
					Name:     name,
					Token:    token,
					ExpToken: expToken,
				}).Times(0)

				return fields{
					service: userMock,
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
			i := UserHandlerImpl{
				userSvc: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.SignUp(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.SignUp(tt.args.w, tt.args.req)
				})
			}

		})
	}
}

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := "Giri Putra Adhittana"
	email := "test@gmail.com"
	password := "123"
	userID := uuid.New()
	token := "123"
	expToken := int64(1000)

	sampleReq := httptest.NewRequest("POST", "http://localhost:8000/v1/sign-in", strings.NewReader(fmt.Sprintf(`{
		"email" : "%s",
		"password" : "%s"
	}`, email, password)))
	sampleResp := httptest.NewRecorder()

	invalidSampleReq := httptest.NewRequest("POST", "http://localhost:8000/v1/sign-in", strings.NewReader(fmt.Sprintf(`{
		"password" : "%s"
	}`, password)))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.UserSvc
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
			name: "success sign in",
			fields: func() fields {
				userMock := mocksvc.NewMockUserSvc(ctrl)

				userMock.EXPECT().SignIn(gomock.Any(), dto.SignInReq{
					Email:    email,
					Password: password,
				}).Return(dto.SignInRes{
					ID:       userID.String(),
					Name:     name,
					Token:    token,
					ExpToken: expToken,
				}).Times(1)

				return fields{
					service: userMock,
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
				userMock := mocksvc.NewMockUserSvc(ctrl)

				userMock.EXPECT().SignIn(gomock.Any(), dto.SignInReq{
					Email:    email,
					Password: password,
				}).Return(dto.SignInRes{
					ID:       userID.String(),
					Name:     name,
					Token:    token,
					ExpToken: expToken,
				}).Times(0)

				return fields{
					service: userMock,
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
			i := UserHandlerImpl{
				userSvc: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.SignIn(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.SignIn(tt.args.w, tt.args.req)
				})
			}

		})
	}
}
