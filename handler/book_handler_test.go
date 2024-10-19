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
	"github.com/gadhittana01/go-modules/utils"
	mockutl "github.com/gadhittana01/go-modules/utils/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBookHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookMock := mocksvc.NewMockBookSvc(ctrl)
	middlewareMock := mockutl.NewMockAuthMiddleware(ctrl)

	type args struct {
		service        service.BookSvc
		authMiddleware utils.AuthMiddleware
	}

	tests := []struct {
		name string
		args args
		want *BookHandlerImpl
	}{
		{
			args: args{
				service:        bookMock,
				authMiddleware: middlewareMock,
			},
			want: &BookHandlerImpl{
				bookSvc:        bookMock,
				authMiddleware: middlewareMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBookHandler(tt.args.service, tt.args.authMiddleware); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBookHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookID := uuid.New()
	title := "Hello"
	description := "World"
	price := float64(100)
	author := "Giri Putra Adhittana"

	sampleReq := httptest.NewRequest("POST", "http://localhost:8000/v1/book", strings.NewReader(fmt.Sprintf(`{
		"title" : "%s",
		"description" : "%s",
		"author" : "%s",
		"price" : %f
	}`, title, description, author, price)))
	sampleResp := httptest.NewRecorder()

	invalidSampleReq := httptest.NewRequest("POST", "http://localhost:8000/v1/book", strings.NewReader(fmt.Sprintf(`{
		"description" : "%s",
		"author" : "%s",
		"price" : %f
	}`, description, author, price)))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.BookSvc
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
			name: "success create book",
			fields: func() fields {
				bookMock := mocksvc.NewMockBookSvc(ctrl)

				bookMock.EXPECT().CreateBook(gomock.Any(), dto.CreateBookReq{
					Title:       title,
					Description: description,
					Author:      author,
					Price:       price,
				}).Return(dto.CreateBookRes{
					ID:          bookID.String(),
					Title:       title,
					Description: description,
					Author:      author,
					Price:       price,
				}).Times(1)

				return fields{
					service: bookMock,
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
				bookMock := mocksvc.NewMockBookSvc(ctrl)

				bookMock.EXPECT().CreateBook(gomock.Any(), dto.CreateBookReq{
					Title:       title,
					Description: description,
					Author:      author,
					Price:       price,
				}).Return(dto.CreateBookRes{
					ID:          bookID.String(),
					Title:       title,
					Description: description,
					Author:      author,
					Price:       price,
				}).Times(0)

				return fields{
					service: bookMock,
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
			i := BookHandlerImpl{
				bookSvc: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.CreateBook(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.CreateBook(tt.args.w, tt.args.req)
				})
			}

		})
	}
}

func TestGetBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookID := uuid.New()
	title := "Hello"
	description := "World"
	price := float64(100)
	author := "Giri Putra Adhittana"
	page := 1
	limit := 10
	totalCount := 20

	sampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/v1/book?page=%d&limit=%d", page, limit), strings.NewReader(``))
	sampleResp := httptest.NewRecorder()

	invalidSampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/v1/book?page=%d&limit=test", page), strings.NewReader(``))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.BookSvc
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
			name: "success get book",
			fields: func() fields {
				bookMock := mocksvc.NewMockBookSvc(ctrl)

				bookMock.EXPECT().GetBook(gomock.Any(), dto.GetBookReq{
					Page:  int32(page),
					Limit: int32(limit),
				}).Return(dto.PaginationResp[dto.GetBookRes]{
					Total:      totalCount,
					IsLoadMore: true,
					Data: []dto.GetBookRes{
						{
							ID:          bookID.String(),
							Title:       title,
							Description: description,
							Author:      author,
							Price:       price,
						},
					},
				}).Times(1)

				return fields{
					service: bookMock,
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
				bookMock := mocksvc.NewMockBookSvc(ctrl)

				bookMock.EXPECT().GetBook(gomock.Any(), dto.GetBookReq{
					Page:  int32(page),
					Limit: int32(limit),
				}).Return(dto.PaginationResp[dto.GetBookRes]{
					Total:      totalCount,
					IsLoadMore: true,
					Data: []dto.GetBookRes{
						{
							ID:          bookID.String(),
							Title:       title,
							Description: description,
							Author:      author,
							Price:       price,
						},
					},
				}).Times(0)

				return fields{
					service: bookMock,
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
			i := BookHandlerImpl{
				bookSvc: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.GetBook(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.GetBook(tt.args.w, tt.args.req)
				})
			}

		})
	}
}
