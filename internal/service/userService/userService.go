package service

import (
	"context"
	"errors"

	"github.com/T-V-N/gophkeeper/internal/config"
	userPB "github.com/T-V-N/gophkeeper/internal/grpc/user"
	"github.com/T-V-N/gophkeeper/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/auth"
	"github.com/T-V-N/gophkeeper/internal/service"
)

type UserService struct {
	userPB.UnimplementedUserServer
	UserApp *app.UserApp
}

func (us *UserService) Register(ctx context.Context, in *userPB.RegisterRequest) (*userPB.RegisterResponse, error) {
	response := userPB.RegisterResponse{}

	uid, err := us.UserApp.Register(ctx, in.Email, in.Password)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(errors.Unwrap(err), &pgErr) {
			return nil, status.Error(codes.AlreadyExists, "already exists")
		}

		switch errors.Unwrap(err) {
		case utils.ErrInvalidEmail:
			return nil, status.Error(codes.InvalidArgument, err.(utils.WrappedAPIError).Message())
		case utils.ErrInvalidPwd:
			return nil, status.Error(codes.InvalidArgument, err.(utils.WrappedAPIError).Message())
		}
	}

	token, err := auth.CreateToken(uid, us.UserApp.Cfg)

	if err != nil {
		return nil, err
	}

	response.AuthToken = token

	return &response, nil
}

func (us *UserService) Login(ctx context.Context, in *userPB.LoginRequest) (*userPB.LoginResponse, error) {
	response := userPB.LoginResponse{}

	token, err := us.UserApp.Login(ctx, in.Email, in.Password, in.OtpCode)

	if err != nil {
		switch errors.Unwrap(err) {
		case utils.ErrBadRequest:
			return nil, status.Error(codes.InvalidArgument, err.(utils.WrappedAPIError).Message())
		case utils.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.(utils.WrappedAPIError).Message())
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, err.(utils.WrappedAPIError).Message())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	response.AuthToken = token

	return &response, nil
}

func (us *UserService) GenerateTOTP(ctx context.Context, _ *userPB.Empty) (*userPB.GenerateTOTPResponse, error) {
	response := userPB.GenerateTOTPResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.(utils.WrappedAPIError).Message())
	}

	otp, err := us.UserApp.GenerateTOTP(ctx, uid)

	if err != nil {
		switch errors.Unwrap(err) {
		case utils.ErrDBLayer:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		case utils.ErrBadRequest:
			return nil, status.Error(codes.InvalidArgument, err.(utils.WrappedAPIError).Message())
		case utils.ErrAppLayer:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	response.OtpKey = otp.Secret()

	return &response, nil
}

func (us *UserService) EnableTOTP(ctx context.Context, in *userPB.EnableTOTPRequest) (*userPB.Empty, error) {
	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.(utils.WrappedAPIError).Message())
	}

	err = us.UserApp.EnableTOTP(ctx, uid, in.OtpCode)

	if err != nil {
		switch errors.Unwrap(err) {
		case utils.ErrDBLayer:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, err.(utils.WrappedAPIError).Message())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	var r userPB.Empty
	return &r, nil
}

func (us *UserService) DisableTOTP(ctx context.Context, in *userPB.DisableTOTPRequest) (*userPB.Empty, error) {
	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.(utils.WrappedAPIError).Message())
	}

	err = us.UserApp.DisableTOTP(ctx, uid, in.OtpCode)

	if err != nil {
		switch errors.Unwrap(err) {
		case utils.ErrDBLayer:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, err.(utils.WrappedAPIError).Message())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var r userPB.Empty
	return &r, nil
}

func InitUserService(cfg *config.Config, a *app.UserApp) *UserService {
	return &UserService{userPB.UnimplementedUserServer{}, a}
}
