package routes

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sheinsviatoslav/gophermart/internal/handlers/addorder"
	"github.com/sheinsviatoslav/gophermart/internal/handlers/addwithdrawal"
	"github.com/sheinsviatoslav/gophermart/internal/handlers/getbalance"
	"github.com/sheinsviatoslav/gophermart/internal/handlers/getorders"
	"github.com/sheinsviatoslav/gophermart/internal/handlers/getwithdrawals"
	"github.com/sheinsviatoslav/gophermart/internal/handlers/login"
	"github.com/sheinsviatoslav/gophermart/internal/handlers/register"
	"github.com/sheinsviatoslav/gophermart/internal/middleware"
	"github.com/sheinsviatoslav/gophermart/internal/storage"
	"log"
	"time"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.WithLogger)
	r.Use(middleware.GzipHandle)
	r.Use(chiMiddleware.Timeout(1000 * time.Millisecond))

	db, err := storage.NewPgStorage()
	if err != nil {
		log.Fatal(err)
	}

	r.Post("/api/user/register", register.NewHandler(db).Handle)
	r.Post("/api/user/login", login.NewHandler(db).Handle)
	r.Get("/api/user/orders", middleware.WithAuth(getorders.NewHandler(db).Handle))
	r.Post("/api/user/orders", middleware.WithAuth(addorder.NewHandler(db).Handle))
	r.Get("/api/user/balance", middleware.WithAuth(getbalance.NewHandler(db).Handle))
	r.Post("/api/user/balance/withdraw", middleware.WithAuth(addwithdrawal.NewHandler(db).Handle))
	r.Get("/api/user/withdrawals", middleware.WithAuth(getwithdrawals.NewHandler(db).Handle))

	return r
}
