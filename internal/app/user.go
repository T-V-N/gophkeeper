package app

import (
	"context"
	"time"

	"github.com/pquerna/otp"
	totp "github.com/pquerna/otp/totp"
	"go.uber.org/zap"

	"github.com/T-V-N/gophkeeper/internal/auth"
	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/helpers"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type EmailSender interface {
	SendConfirmationEmail(to, confirmationURL string) error
}

type User interface {
	CreateUser(ctx context.Context, email, passwordHash, confirmationCode string) (string, error)
	GetUserByEmail(ctx context.Context, email string) (storage.User, error)
	GetUserByID(ctx context.Context, uid string) (storage.User, error)
	UpdateUser(ctx context.Context, uid, email, passwordHash, totpSecret string, totpEnabled bool, confirmedAt time.Time) error
	Close()
}

type UserApp struct {
	User        User
	Cfg         *config.Config
	EmailSender EmailSender
	logger      *zap.SugaredLogger
}

func InitUserApp(cfg *config.Config, logger *zap.SugaredLogger, emailSender *helpers.EmailSender) (*UserApp, error) {
	user, err := storage.InitUserStorage(cfg)

	if err != nil {
		return nil, err
	}

	return &UserApp{user, cfg, emailSender, logger}, nil
}

func (app *UserApp) Register(ctx context.Context, email, password string) (string, error) {
	if !utils.IsValidEmail(email) {
		return "", utils.ErrInvalidEmail
	}

	if !utils.IsValidPassword(password) {
		return "", utils.ErrInvalidPwd
	}

	passwordHash, err := utils.HashDataSecurely(password)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrAppLayer)
	}

	confirmationCode := utils.GenerateConfirmationCode(email, app.Cfg.SecretKey)

	uid, err := app.User.CreateUser(ctx, email, passwordHash, confirmationCode)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrDBLayer)
	}

	err = app.EmailSender.SendConfirmationEmail(email, confirmationCode)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrThirdParty)
	}

	return uid, nil
}

func (app *UserApp) ConfirmUser(ctx context.Context, email, code string) error {
	user, err := app.User.GetUserByEmail(ctx, email)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	if user.VerificationCode != code {
		return utils.WrapError(utils.ErrAuth, nil)
	}

	err = app.User.UpdateUser(ctx, user.UID, user.Email, user.PasswordHash, user.TOTPSecret, user.TOTPEnabled, time.Now())

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	return nil
}

func (app *UserApp) Login(ctx context.Context, email, password, otpCode string) (string, error) {
	user, err := app.User.GetUserByEmail(ctx, email)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrNotFound)
	}

	if user.TOTPEnabled && !totp.Validate(otpCode, user.TOTPSecret) {
		return "", utils.WrapError(utils.ErrBadRequest, &utils.APIError{Msg: "Invalid OTP code"})
	}

	isPasswordValid := utils.CheckPasswordHash(password, user.PasswordHash)

	if !isPasswordValid {
		return "", utils.WrapError(utils.ErrNotAuthorized, &utils.APIError{Msg: "Invalid password"})
	}

	token, err := auth.CreateToken(user.UID, app.Cfg)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrAppLayer)
	}

	return token, nil
}

func (app *UserApp) GenerateTOTP(ctx context.Context, uid string) (*otp.Key, error) {
	user, err := app.User.GetUserByID(ctx, uid)

	if err != nil {
		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	if user.TOTPEnabled {
		return nil, utils.WrapError(utils.ErrBadRequest, &utils.APIError{Msg: "TOTP already enabled"})
	}

	key, err := totp.Generate(totp.GenerateOpts{
		AccountName: uid,
		Issuer:      "gophkeeper",
	})

	if err != nil {
		return nil, utils.WrapError(err, utils.ErrAppLayer)
	}

	err = app.User.UpdateUser(ctx, user.UID, user.Email, user.PasswordHash, key.Secret(), user.TOTPEnabled, user.ConfirmedAt)

	if err != nil {
		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	return key, nil
}

func (app *UserApp) EnableTOTP(ctx context.Context, uid, otpCode string) error {
	user, err := app.User.GetUserByID(ctx, uid)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	if user.TOTPEnabled {
		return nil
	}

	if !totp.Validate(otpCode, user.TOTPSecret) {
		return utils.WrapError(utils.ErrNotAuthorized, &utils.APIError{Msg: "Wrong TOTP code"})
	}

	err = app.User.UpdateUser(ctx, user.UID, user.Email, user.PasswordHash, user.TOTPSecret, true, user.ConfirmedAt)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	return nil
}

func (app *UserApp) DisableTOTP(ctx context.Context, uid, otpCode string) error {
	user, err := app.User.GetUserByID(ctx, uid)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	if !user.TOTPEnabled {
		return nil
	}

	if !totp.Validate(otpCode, user.TOTPSecret) {
		return utils.WrapError(utils.ErrNotAuthorized, &utils.APIError{Msg: "Wrong TOTP code"})
	}

	err = app.User.UpdateUser(ctx, user.UID, user.Email, user.PasswordHash, "", false, user.ConfirmedAt)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	return nil
}
