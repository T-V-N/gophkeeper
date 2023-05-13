package app_test

import (
	"context"
	"testing"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/helpers"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_e2e_register_upload_everything(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	sugar := logger.Sugar()

	cfg, err := config.Init()
	if err != nil {
		sugar.Fatalw("Unable to load config",
			"Error", err,
		)
	}

	t.Run("Connects to storage", func(t *testing.T) {
		_, err := storage.InitStorage(*cfg)

		assert.NoError(t, err, "Shall connect")
	})

	s, err := storage.InitStorage(*cfg)
	sender := helpers.InitEmailSender(cfg)
	userApp := &app.UserApp{}
	t.Run("Creates user", func(t *testing.T) {
		userApp, err = app.InitUserApp(s.Conn, cfg, sugar, *sender)
		uid, err := userApp.Register(context.Background(), "dr.tvn@yandex.ru", "password1!A")
		user, err := userApp.User.GetUserByID(context.Background(), uid)
		err = userApp.ConfirmUser(context.Background(), "dr.tvn@yandex.ru", user.VerificationCode)

		assert.NoError(t, err, "Shall work without any probelem")
	})

	// cardApp, err := app.InitCardApp(s.Conn, cfg, sugar)
	// fileApp, err := app.InitFileApp(s.Conn, cfg, sugar)
	// passwordApp, err := app.InitLogPasswordApp(s.Conn, cfg, sugar)
	// noteApp, err := app.InitTextNoteApp(s.Conn, cfg, sugar)

}
