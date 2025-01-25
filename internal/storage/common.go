package storage

import (
	"context"
	"github.com/sheinsviatoslav/gophermart/internal/common"
)

type Storage interface {
	CreateUser(context.Context, common.UserCredentials) (common.User, error)
	GetUserByLogin(context.Context, string) (common.User, error)
	CheckLoginExists(context.Context, string) (bool, error)
	AddOrder(context.Context, string, string) error
	GetOrderByNumber(context.Context, string) (common.Order, error)
	GetUserOrders(context.Context, string) ([]common.Order, error)
	GetUserBalance(context.Context, string) (common.Balance, error)
	AddWithdrawal(context.Context, float64, string, string) error
	GetUserWithdrawals(context.Context, string) ([]common.Withdrawal, error)
	GetOrderFromAccrual(string) (common.Order, error)
	UpdateOrderFromAccrual(context.Context, chan struct{}, string, string) error
}
