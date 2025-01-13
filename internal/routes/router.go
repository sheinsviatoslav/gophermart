package routes

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
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
	//r.Use(middleware.GzipHandle)
	r.Use(chiMiddleware.Timeout(1000 * time.Millisecond))

	db := storage.NewPgStorage()
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	r.Post("/api/user/register", register.NewHandler(db).Handle)
	r.Post("/api/user/login", login.NewHandler(db).Handle)

	return r
}
