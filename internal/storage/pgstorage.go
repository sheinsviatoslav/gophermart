package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sheinsviatoslav/gophermart/internal/common"
	"github.com/sheinsviatoslav/gophermart/internal/config"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

type PgStorage struct {
	DB *sql.DB
}

func NewPgStorage() (*PgStorage, error) {
	p := &PgStorage{
		DB: nil,
	}

	var err error
	p.DB, err = sql.Open("pgx", *config.DatabaseURI)
	if err != nil {
		return nil, err
	}

	_, err = p.DB.Exec(
		"CREATE TABLE IF NOT EXISTS users (" +
			"id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(), " +
			"login TEXT NOT NULL UNIQUE, " +
			"password TEXT NOT NULL," +
			"current NUMERIC(10, 2) DEFAULT 0," +
			"withdrawn NUMERIC(10, 2) DEFAULT 0);" +
			"CREATE TABLE IF NOT EXISTS orders (" +
			"id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(), " +
			"number TEXT NOT NULL UNIQUE," +
			"status TEXT DEFAULT 'NEW'," +
			"accrual NUMERIC(10, 2) DEFAULT 0," +
			"user_id UUID REFERENCES users(id)," +
			"uploaded_at TIMESTAMP NOT NULL DEFAULT now());" +
			"CREATE TABLE IF NOT EXISTS withdrawals (" +
			"id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(), " +
			"sum NUMERIC(10, 2) NOT NULL," +
			"\"order\" TEXT UNIQUE," +
			"processed_at TIMESTAMP NOT NULL DEFAULT now(), " +
			"user_id UUID REFERENCES users(id));")
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *PgStorage) CreateUser(ctx context.Context, user common.UserCredentials) (common.User, error) {
	newUser := common.User{}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return newUser, err
	}

	tx, err := p.DB.Begin()
	if err != nil {
		return newUser, err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (login, password) VALUES($1, $2)`
	if _, err = tx.ExecContext(ctx, query, user.Login, hashedPassword); err != nil {
		return newUser, err
	}

	if err = tx.Commit(); err != nil {
		return newUser, err
	}

	newUser, err = p.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (p *PgStorage) GetUserByLogin(ctx context.Context, login string) (common.User, error) {
	newUser := common.User{}
	query := `SELECT id, login, password FROM users WHERE login = $1`
	row := p.DB.QueryRowContext(ctx, query, login)
	if err := row.Scan(&newUser.ID, &newUser.Login, &newUser.Password); err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (p *PgStorage) CheckLoginExists(ctx context.Context, login string) (bool, error) {
	var exists *bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1) AS exists`
	if err := p.DB.QueryRowContext(ctx, query, login).Scan(&exists); err != nil {
		return *exists, err
	}

	return *exists, nil
}

func (p *PgStorage) AddOrder(ctx context.Context, orderNumber string, userID string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO orders (number, user_id) VALUES($1, $2)`
	if _, err = tx.ExecContext(ctx, query, orderNumber, userID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) GetOrderByNumber(ctx context.Context, orderNumber string) (common.Order, error) {
	var order common.Order
	query := `SELECT number, user_id FROM orders WHERE number = $1`
	row := p.DB.QueryRowContext(ctx, query, orderNumber)

	err := row.Scan(&order.Number, &order.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return order, nil
	}

	if err != nil {
		return order, err
	}

	return order, nil
}

func (p *PgStorage) GetUserOrders(ctx context.Context, userID string) ([]common.Order, error) {
	var orders []common.Order
	query := `SELECT id, number, uploaded_at, user_id, accrual, status FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := p.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var order common.Order

		if err = rows.Scan(
			&order.ID,
			&order.Number,
			&order.UploadedAt,
			&order.UserID,
			&order.Accrual,
			&order.Status,
		); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (p *PgStorage) GetOrderFromAccrual(orderNumber string) (common.Order, error) {
	var order common.Order
	response, err := http.Get(*config.AccrualSystemAddress + "/api/orders/" + orderNumber)
	if err != nil {
		return order, err
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return order, err
	}

	if err = json.Unmarshal(body, &order); err != nil {
		return order, err
	}

	return order, nil
}

func (p *PgStorage) UpdateOrderFromAccrual(ctx context.Context, done chan struct{}, orderNumber string, userID string) error {
	order, err := p.GetOrderFromAccrual(orderNumber)
	if err != nil {
		return err
	}

	if order.Status == "INVALID" || order.Status == "PROCESSED" {
		close(done)
	}

	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`
	if _, err = tx.ExecContext(ctx, query, order.Status, order.Accrual, orderNumber); err != nil {
		return err
	}

	query = `UPDATE users SET current = current + $1 WHERE id = $2`
	if _, err = tx.ExecContext(ctx, query, order.Accrual, userID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) GetUserBalance(ctx context.Context, userID string) (common.Balance, error) {
	var balance common.Balance
	query := `SELECT current, withdrawn FROM users WHERE id = $1`
	row := p.DB.QueryRowContext(ctx, query, userID)

	err := row.Scan(&balance.Current, &balance.Withdrawn)
	if errors.Is(err, sql.ErrNoRows) {
		return balance, nil
	}

	if err != nil {
		return balance, err
	}

	return balance, nil
}

func (p *PgStorage) AddWithdrawal(ctx context.Context, sum float64, orderNumber, userID string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO withdrawals (sum, "order", user_id) VALUES($1, $2, $3)`
	if _, err = tx.ExecContext(ctx, query, sum, orderNumber, userID); err != nil {
		return err
	}

	query = `UPDATE users SET current = current - $1, withdrawn = withdrawn + $1 WHERE id = $2`
	if _, err = tx.ExecContext(ctx, query, sum, userID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) GetUserWithdrawals(ctx context.Context, userID string) ([]common.Withdrawal, error) {
	var withdrawals []common.Withdrawal
	query := `SELECT id, sum, "order", processed_at FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC`
	rows, err := p.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdrawal common.Withdrawal

		if err = rows.Scan(&withdrawal.ID, &withdrawal.Sum, &withdrawal.Order, &withdrawal.ProcessedAt); err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
