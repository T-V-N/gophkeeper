package app

import (
	"context"
	"time"

	"github.com/T-V-N/gophkeeper/internal/auth"
	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/helpers"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pquerna/otp"
	totp "github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

type EmailSender interface {
	SendConfirmationEmail(to, confirmationURL string) error
}

type UserApp struct {
	User        *storage.UserStorage
	Cfg         *config.Config
	logger      *zap.SugaredLogger
	emailSender EmailSender
}

func InitUserApp(conn *pgxpool.Pool, cfg *config.Config, logger *zap.SugaredLogger, emailSender helpers.EmailSender) (*UserApp, error) {
	user, err := storage.InitUser(conn)

	if err != nil {
		return nil, err
	}

	return &UserApp{user, cfg, logger, emailSender}, nil
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

	confirmationCode := utils.String(10)

	err = app.emailSender.SendConfirmationEmail(email, confirmationCode)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrThirdParty)
	}

	uid, err := app.User.CreateUser(ctx, email, passwordHash, confirmationCode)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrDBLayer)
	}

	return uid, nil
}

func (app *UserApp) ConfirmUser(ctx context.Context, email, code string) error {
	user, err := app.User.GetUserByEmail(ctx, email)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	if user.VerificationCode != code {
		return utils.ErrAuth
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
		return "", utils.WrapError(err, utils.ErrDBLayer)
	}

	if user.TOTPEnabled && !totp.Validate(otpCode, user.TOTPSecret) {
		return "", utils.ErrNotAuthorized
	}

	if err != nil {
		return "", utils.ErrNotAuthorized
	}

	isPasswordValid := utils.CheckPasswordHash(password, user.PasswordHash)

	if !isPasswordValid {
		return "", utils.ErrNotAuthorized
	}

	return auth.CreateToken(user.UID, app.Cfg)
}

func (app *UserApp) GenerateTOTP(ctx context.Context, uid string) (*otp.Key, error) {
	user, err := app.User.GetUserByID(ctx, uid)

	if err != nil {
		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	if user.TOTPEnabled {
		return nil, utils.ErrBadRequest
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
		return utils.ErrNotAuthorized
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
		return utils.ErrNotAuthorized
	}

	err = app.User.UpdateUser(ctx, user.UID, user.Email, user.PasswordHash, "", false, user.ConfirmedAt)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	return nil
}
