package storage

import (
	"context"
	"github.com/sheinsviatoslav/gophermart/internal/common"
)

type Storage interface {
	CreateUser(context.Context, common.User) error
	GetUserPasswordByLogin(context.Context, string) (string, error)
	CheckLoginExists(context.Context, string) (bool, error)
}
