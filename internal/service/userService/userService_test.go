package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/caarlos0/env/v8"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	userPB "github.com/T-V-N/gophkeeper/internal/grpc/user"
	mock "github.com/T-V-N/gophkeeper/internal/mocks"
	service "github.com/T-V-N/gophkeeper/internal/service/userService"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

func TestUserService_Register(t *testing.T) {
	type want struct {
		authToken string
		errCode   codes.Code
	}

	type StorageResponse struct {
		id  string
		err error
	}

	tests := []struct {
		name            string
		request         userPB.RegisterRequest
		want            want
		storageResponse StorageResponse
		sendEmail       bool
		accessStore     bool
		senderResponse  error
	}{
		{
			name: "register normally",
			request: userPB.RegisterRequest{
				Email:    "testmail@yandex.ru",
				Password: "testpassword1A!",
			},
			want: want{
				errCode: codes.OK,
			},
			accessStore:     true,
			storageResponse: StorageResponse{"1", nil},
			sendEmail:       true,
			senderResponse:  nil,
		},
		{
			name: "duplicate",
			request: userPB.RegisterRequest{
				Email:    "testmail@yandex.ru",
				Password: "testpassword1A!",
			},
			want: want{
				errCode: codes.AlreadyExists,
			},
			accessStore:     true,
			storageResponse: StorageResponse{"", utils.ErrDuplicate},
			sendEmail:       false,
			senderResponse:  nil,
		},
		{
			name: "wrong email",
			request: userPB.RegisterRequest{
				Email:    "dfjdsnfddnfijdwfidsjfidsjfisdjf",
				Password: "testpassword1A!",
			},
			want: want{
				errCode: codes.InvalidArgument,
			},
			accessStore:    false,
			sendEmail:      false,
			senderResponse: nil,
		},
		{
			name: "wrong password",
			request: userPB.RegisterRequest{
				Email:    "testmail@yandex.ru",
				Password: "",
			},
			want: want{
				errCode: codes.InvalidArgument,
			},
			accessStore:     false,
			storageResponse: StorageResponse{"", &pgconn.PgError{Code: pgerrcode.UniqueViolation}},
			sendEmail:       false,
			senderResponse:  nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)
	user := mock.NewMockUser(ctrl)
	sender := mock.NewMockEmailSender(ctrl)
	ua := app.UserApp{User: user, Cfg: cfg, EmailSender: sender}
	us := service.UserService{userPB.UnimplementedUserServer{}, &ua}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confirmationCode := utils.GenerateConfirmationCode(tt.request.Email, cfg.SecretKey)
			if tt.accessStore {
				user.EXPECT().CreateUser(context.Background(), tt.request.Email, gomock.Any(), confirmationCode).Return(tt.storageResponse.id, tt.storageResponse.err)
			}
			if tt.sendEmail {
				sender.EXPECT().SendConfirmationEmail(tt.request.Email, gomock.Any()).Return(tt.senderResponse)
			}

			_, err := us.Register(context.Background(), &tt.request)

			if err != nil {
				if unwrapError, ok := status.FromError(err); ok {
					assert.Equal(t, unwrapError.Code(), tt.want.errCode)
				} else {
					t.Error("can't parse error, unexpected test result")
				}
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)
	user := mock.NewMockUser(ctrl)
	ua := app.UserApp{User: user, Cfg: cfg}
	us := service.UserService{userPB.UnimplementedUserServer{}, &ua}

	t.Run("Login normally, no otp", func(t *testing.T) {
		request := userPB.LoginRequest{
			Email:    "hey@gmail.com",
			Password: "123123123",
		}

		passwordHash, err := utils.HashDataSecurely(request.Password)

		user.EXPECT().GetUserByEmail(context.Background(), request.Email).Return(&storage.User{
			UID:          "1",
			TOTPEnabled:  false,
			PasswordHash: passwordHash,
		}, nil)

		_, err = us.Login(context.Background(), &request)

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.OK)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})
	t.Run("Login normally, with otp", func(t *testing.T) {
		secret := "AAAAAAAA"
		otpCode, err := totp.GenerateCode(secret, time.Now())

		request := userPB.LoginRequest{
			Email:    "hey@gmail.com",
			Password: "123123123",
			OtpCode:  otpCode,
		}

		passwordHash, _ := utils.HashDataSecurely(request.Password)
		user.EXPECT().GetUserByEmail(context.Background(), request.Email).Return(&storage.User{
			UID:          "1",
			PasswordHash: passwordHash,
			TOTPSecret:   "AAAAAAAA",
			TOTPEnabled:  true,
		}, nil)

		_, err = us.Login(context.Background(), &request)

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.OK)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})
	t.Run("Login normally, but wrong otp", func(t *testing.T) {
		secret := "AABAAAAA"
		otpCode, err := totp.GenerateCode(secret, time.Now())

		request := userPB.LoginRequest{
			Email:    "hey@gmail.com",
			Password: "123123123",
			OtpCode:  otpCode,
		}

		passwordHash, _ := utils.HashDataSecurely(request.Password)
		user.EXPECT().GetUserByEmail(context.Background(), request.Email).Return(&storage.User{
			UID:          "1",
			PasswordHash: passwordHash,
			TOTPSecret:   "AAAAAAAA",
			TOTPEnabled:  true,
		}, nil)

		_, err = us.Login(context.Background(), &request)

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.InvalidArgument)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})

	t.Run("Login with wrong email", func(t *testing.T) {
		request := userPB.LoginRequest{
			Email:    "hey@gmail.com",
			Password: "123123123",
		}

		user.EXPECT().GetUserByEmail(context.Background(), request.Email).Return(nil, utils.ErrNotFound)

		_, err := us.Login(context.Background(), &request)

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.NotFound)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})

	t.Run("Login with wrong password", func(t *testing.T) {
		request := userPB.LoginRequest{
			Email:    "hey@gmail.com",
			Password: "123123123",
		}

		passwordHash, _ := utils.HashDataSecurely(request.Password + "1111")
		user.EXPECT().GetUserByEmail(context.Background(), request.Email).Return(&storage.User{
			UID:          "1",
			PasswordHash: passwordHash,
		}, nil)

		_, err := us.Login(context.Background(), &request)

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.Unauthenticated)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})
}

func TestUserService_OTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)
	user := mock.NewMockUser(ctrl)
	ua := app.UserApp{User: user, Cfg: cfg}
	us := service.UserService{userPB.UnimplementedUserServer{}, &ua}

	t.Run("Generate OTP Secret", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1"))

		user.EXPECT().GetUserByID(ctx, "1").Return(&storage.User{
			UID:         "1",
			TOTPEnabled: false,
		}, nil)

		user.EXPECT().UpdateUser(ctx, "1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		key, err := us.GenerateTOTP(ctx, &userPB.Empty{})

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.OK)
			assert.NotEmpty(t, key)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})

	t.Run("Generate OTP Secret when already enabled", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1"))

		user.EXPECT().GetUserByID(ctx, "1").Return(&storage.User{
			UID:         "1",
			TOTPEnabled: true,
		}, nil)

		key, err := us.GenerateTOTP(ctx, &userPB.Empty{})

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.InvalidArgument)
			assert.Empty(t, key)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})

	t.Run("Confirm secret with wrong code", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1"))

		user.EXPECT().GetUserByID(ctx, "1").Return(&storage.User{
			UID:         "1",
			TOTPEnabled: false,
		}, nil)

		_, err := us.EnableTOTP(ctx, &userPB.EnableTOTPRequest{
			OtpCode: "this code wrong",
		})

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.InvalidArgument)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})

	t.Run("Confirm secret with good code", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1"))
		otpKey, err := totp.Generate(totp.GenerateOpts{
			AccountName: "1",
			Issuer:      "gophkeeper",
		})

		otpCode, err := totp.GenerateCode(otpKey.Secret(), time.Now())

		user.EXPECT().GetUserByID(ctx, "1").Return(&storage.User{
			UID:         "1",
			TOTPEnabled: false,
			TOTPSecret:  otpKey.Secret(),
		}, nil)
		user.EXPECT().UpdateUser(ctx, "1", gomock.Any(), gomock.Any(), otpKey.Secret(), gomock.Any(), gomock.Any()).Return(nil)

		_, err = us.EnableTOTP(ctx, &userPB.EnableTOTPRequest{
			OtpCode: otpCode,
		})

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.OK)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})

	t.Run("Confirm secret when it's enabled already", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1"))

		user.EXPECT().GetUserByID(ctx, "1").Return(&storage.User{
			UID:         "1",
			TOTPEnabled: true,
		}, nil)

		_, err := us.EnableTOTP(ctx, &userPB.EnableTOTPRequest{
			OtpCode: "any code",
		})

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.OK)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})

	t.Run("Disable secret with wrong code", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1"))
		otpKey, err := totp.Generate(totp.GenerateOpts{
			AccountName: "1",
			Issuer:      "gophkeeper",
		})

		user.EXPECT().GetUserByID(ctx, "1").Return(&storage.User{
			UID:         "1",
			TOTPEnabled: true,
			TOTPSecret:  otpKey.Secret(),
		}, nil)

		_, err = us.DisableTOTP(ctx, &userPB.DisableTOTPRequest{
			OtpCode: "code wrong",
		})

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.Unauthenticated)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})

	t.Run("Disable secret with good code", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1"))
		otpKey, err := totp.Generate(totp.GenerateOpts{
			AccountName: "1",
			Issuer:      "gophkeeper",
		})

		otpCode, err := totp.GenerateCode(otpKey.Secret(), time.Now())
		user.EXPECT().GetUserByID(ctx, "1").Return(&storage.User{
			UID:         "1",
			TOTPEnabled: true,
			TOTPSecret:  otpKey.Secret(),
		}, nil)

		user.EXPECT().UpdateUser(ctx, "1", gomock.Any(), gomock.Any(), "", false, gomock.Any()).Return(nil)
		_, err = us.DisableTOTP(ctx, &userPB.DisableTOTPRequest{
			OtpCode: otpCode,
		})

		if unwrapError, ok := status.FromError(err); ok {
			assert.Equal(t, unwrapError.Code(), codes.OK)
		} else {
			t.Error("can't parse error, unexpected test result")
		}
	})
}
