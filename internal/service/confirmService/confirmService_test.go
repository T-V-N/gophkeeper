package service_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caarlos0/env/v8"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	mock "github.com/T-V-N/gophkeeper/internal/mocks"
	service "github.com/T-V-N/gophkeeper/internal/service/confirmService"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

func TestConfirmService_Handle(t *testing.T) {
	type want struct {
		authToken string
		errCode   codes.Code
	}

	type StorageResponse struct {
		id  string
		err error
	}

	tests := []struct {
		name               string
		want               want
		email              string
		code               string
		getByEmailResponse storage.User
		getByEmailError    error
		updateUser         bool
	}{
		{
			name: "confirm normally",
			want: want{
				errCode: http.StatusOK,
			},
			email:              "test@test.com",
			code:               "hey123",
			getByEmailResponse: storage.User{UID: "1", Email: "test@test.com", VerificationCode: "hey123"},
			updateUser:         true,
		},
		{
			name: "confirm with wrong code",
			want: want{
				errCode: http.StatusBadRequest,
			},
			email:              "test@test.com",
			code:               "WRONG!!!!!!!!",
			getByEmailResponse: storage.User{UID: "1", Email: "test@test.com", VerificationCode: "hey123"},
			updateUser:         false,
		},
		{
			name: "confirm with no code",
			want: want{
				errCode: http.StatusBadRequest,
			},
			email:              "test@test.com",
			code:               "",
			getByEmailResponse: storage.User{UID: "1", Email: "test@test.com", VerificationCode: "hey123"},
			updateUser:         false,
		},
		{
			name: "confirm with no email and code",
			want: want{
				errCode: http.StatusBadRequest,
			},
			email:              "",
			code:               "",
			getByEmailResponse: storage.User{UID: "1", Email: "test@test.com", VerificationCode: "hey123"},
			updateUser:         false,
		},
		{
			name: "confirm user which doesn't exist",
			want: want{
				errCode: http.StatusNoContent,
			},
			email:              "non_existent@email.com",
			code:               "11111111",
			getByEmailResponse: storage.User{},
			getByEmailError:    utils.ErrNotFound,
			updateUser:         false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	user := mock.NewMockUser(ctrl)

	sender := mock.NewMockEmailSender(ctrl)
	ua := app.UserApp{User: user, Cfg: cfg, EmailSender: sender}

	httpHandler := service.InitConfirmationService(&ua)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.email != "" && tt.code != "" {
				user.EXPECT().GetUserByEmail(gomock.Any(), tt.email).Return(&tt.getByEmailResponse, tt.getByEmailError)
			}
			if tt.updateUser {
				user.EXPECT().UpdateUser(gomock.Any(), "1", tt.email, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			}

			request := httptest.NewRequest(http.MethodGet, "/confirm", nil)
			w := httptest.NewRecorder()
			ctx := chi.NewRouteContext()

			ctx.URLParams.Add("confirmationCode", tt.code)
			ctx.URLParams.Add("email", tt.email)

			rctx := context.WithValue(request.Context(), chi.RouteCtxKey, ctx)
			request = request.WithContext(rctx)
			httpHandler.HandleConfirmUser(w, request)

			res := w.Result()
			res.Body.Close()

			assert.Equal(t, res.StatusCode, int(tt.want.errCode))
		})
	}
}
