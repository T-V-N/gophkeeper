package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/utils"
	"github.com/stretchr/testify/assert"

	"github.com/T-V-N/gophkeeper/internal/storage"
)

func InitTestConfig() config.Config {
	return config.Config{DatabaseURI: "postgresql://root:root@localhost:5433/keeper"}
}

func Test_InitStorage(t *testing.T) {
	cfg := InitTestConfig()

	t.Run("Connects to storage", func(t *testing.T) {
		_, err := storage.InitStorage(cfg)

		assert.NoError(t, err, "Shall connect")
	})
}

func Test_User(t *testing.T) {
	cfg := InitTestConfig()
	st, _ := storage.InitStorage(cfg)
	user := storage.UserStorage{Conn: st.Conn}

	var id string

	t.Run("Create user", func(t *testing.T) {
		uid, err := user.CreateUser(context.Background(), "hey@gmail.com", "xxx", "kekw")
		id = uid

		assert.NoError(t, err, "Shall create")
		assert.NotEmpty(t, uid, "ID shall exist")
	})

	t.Run("Create user but duplicate", func(t *testing.T) {
		_, err := user.CreateUser(context.Background(), "hey@gmail.com", "xxx", "kekw")

		assert.Error(t, err, utils.ErrDuplicate)
	})

	t.Run("Update user", func(t *testing.T) {
		err := user.UpdateUser(context.Background(), id, "new@email.com", "xxx", "xxxx", true, time.Now())

		assert.NoError(t, err, utils.ErrDuplicate)
	})

	t.Run("Find user by id", func(t *testing.T) {

	})

	t.Run("Find user by email", func(t *testing.T) {

	})
}

func Test_Password(t *testing.T) {

}

func Test_File(t *testing.T) {
}
