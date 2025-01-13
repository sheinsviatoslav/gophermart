package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sheinsviatoslav/gophermart/internal/common"
	"github.com/sheinsviatoslav/gophermart/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type PgStorage struct {
	DB *sql.DB
}

func NewPgStorage() *PgStorage {
	return &PgStorage{
		DB: nil,
	}
}

func (p *PgStorage) Connect() error {
	var err error
	p.DB, err = sql.Open("pgx", *config.DatabaseURI)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec(
		"CREATE TABLE IF NOT EXISTS users (" +
			"id uuid PRIMARY KEY DEFAULT GEN_RANDOM_UUID(), " +
			"login TEXT NOT NULL UNIQUE, " +
			"password TEXT NOT NULL)")
	if err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) CreateUser(ctx context.Context, user common.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (login, password) VALUES($1, $2)`
	if _, err = tx.ExecContext(ctx, query, user.Login, hashedPassword); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) GetUserPasswordByLogin(ctx context.Context, login string) (string, error) {
	var password *string
	query := `SELECT password from users WHERE login = $1`
	row := p.DB.QueryRowContext(ctx, query, login)
	if err := row.Scan(&password); err != nil {
		return "", err
	}

	return *password, nil
}

func (p *PgStorage) CheckLoginExists(ctx context.Context, login string) (bool, error) {
	var exists *bool
	query := `SELECT EXISTS(SELECT 1 from users WHERE login = $1) AS exists`
	if err := p.DB.QueryRowContext(ctx, query, login).Scan(&exists); err != nil {
		return *exists, err
	}

	return *exists, nil
}
