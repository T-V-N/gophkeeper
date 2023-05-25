package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/auth"
	"github.com/T-V-N/gophkeeper/internal/config"
	userPB "github.com/T-V-N/gophkeeper/internal/grpc/user"
	"github.com/T-V-N/gophkeeper/internal/service"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type UserService struct {
	userPB.UnimplementedUserServer
	UserApp *app.UserApp
}

func (us *UserService) Register(ctx context.Context, in *userPB.RegisterRequest) (*userPB.RegisterResponse, error) {
	response := userPB.RegisterResponse{}

	uid, err := us.UserApp.Register(ctx, in.Email, in.Password)

	if err != nil {
		switch err {
		case utils.ErrInvalidEmail:
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		case utils.ErrInvalidPwd:
			return nil, status.Error(codes.InvalidArgument, "invalid password")
		case utils.ErrDuplicate:
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}

	token, err := auth.CreateToken(uid, us.UserApp.Cfg)

	if err != nil {
		return nil, status.Error(codes.Internal, "cannot create token")
	}

	response.AuthToken = token

	return &response, nil
}

func (us *UserService) Login(ctx context.Context, in *userPB.LoginRequest) (*userPB.LoginResponse, error) {
	response := userPB.LoginResponse{}

	token, err := us.UserApp.Login(ctx, in.Email, in.Password, in.OtpCode)

	if err != nil {
		switch err {
		case utils.ErrInvalidTOTP:
			return nil, status.Error(codes.InvalidArgument, "invalid totp")
		case utils.ErrInvalidEmail:
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		case utils.ErrInvalidPwd:
			return nil, status.Error(codes.Unauthenticated, "wrong password")
		case utils.ErrNotFound:
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}

	response.AuthToken = token

	return &response, nil
}

func (us *UserService) GenerateTOTP(ctx context.Context, _ *userPB.Empty) (*userPB.GenerateTOTPResponse, error) {
	response := userPB.GenerateTOTPResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	otp, err := us.UserApp.GenerateTOTP(ctx, uid)

	if err != nil {
		switch err {
		case utils.ErrBadRequest:
			return nil, status.Error(codes.InvalidArgument, "totp already enabled")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}

	response.OtpKey = otp.Secret()

	return &response, nil
}

func (us *UserService) EnableTOTP(ctx context.Context, in *userPB.EnableTOTPRequest) (*userPB.Empty, error) {
	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = us.UserApp.EnableTOTP(ctx, uid, in.OtpCode)

	if err != nil {
		switch err {
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.InvalidArgument, "wrong otp code")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}
	var r userPB.Empty
	return &r, nil
}

func (us *UserService) DisableTOTP(ctx context.Context, in *userPB.DisableTOTPRequest) (*userPB.Empty, error) {
	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = us.UserApp.DisableTOTP(ctx, uid, in.OtpCode)

	if err != nil {
		switch err {
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, "wrong otp code")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}

	var r userPB.Empty
	return &r, nil
}

func InitUserService(cfg *config.Config, a *app.UserApp) *UserService {
	return &UserService{userPB.UnimplementedUserServer{}, a}
}
