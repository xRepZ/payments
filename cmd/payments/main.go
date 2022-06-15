package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	logger "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xRepZ/payments/cmd/config"

	api "github.com/xRepZ/payments/internal/http"
	postgres "github.com/xRepZ/payments/internal/storage"
)

func main() {
	log := logger.New()
	ctx, cancel := context.WithCancel(context.Background())

	// парсим конфиг
	var cfgPath string
	flag.StringVar(&cfgPath, "c", "./config/config.yaml", "path to config")
	flag.Parse()

	cfg, err := config.ParseConfig(cfgPath)
	if err != nil {
		log.WithError(err).Fatal("can't parse config")
	}

	defer func() {
		msg := recover()
		if msg != nil {
			log.WithField("panic_msg", msg).Error("recover from panic and exit")
			cancel()
		}
	}()
	// ждём базу из докер компоуза
	time.Sleep(time.Second * 10)
	db, err := sqlx.Connect("postgres", cfg.Storage.Postgres.Dsn)
	if err != nil {
		log.WithError(err).Fatal("can't connect to postgres")
	}

	storage := postgres.NewTransaction(db)

	server := api.NewServer(storage)

	// Объявляем handlers
	router := chi.NewRouter()

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(api.TokenGenerator))
		r.Use(jwtauth.Authenticator)

		r.Patch("/api/transaction/{transaction_id}", server.UpdateById) // Изменение статуса платежа

	})

	// public
	router.Group(func(r chi.Router) {
		r.Get("/api/token", server.Login)
		// Вызывается пользователем
		r.Post("/api/transaction", server.AddTransaction)                  // Создание платежа
		r.Get("/api/transaction/{transaction_id}", server.GetById)         // Проверка статуса платежа по id
		r.Get("/api/transaction/user/{user_id}", server.GetByUserId)       // Все платежи юзера по его id
		r.Get("/api/transaction/user/mail/{user_email}", server.GetByMail) // Все платежи юзера по его email
		r.Delete("/api/transaction/{transaction_id}", server.CancelById)   // Отмена платежа по его id
		// Вызывается платежной системой
		//r.With(AuthMiddleware).Patch("/api/transaction/{transaction_id}", ) // Изменение статуса платежа
	})

	wg := sync.WaitGroup{}

	wg.Add(1)
	osSigCh := make(chan os.Signal, 1)
	signal.Notify(osSigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer wg.Done()

		// Программа будет завершена либо из-за контекста, либо по сигналу
		select {
		case <-ctx.Done():
		case s := <-osSigCh:
			log.Info("shutdown by singal: ", s.String())
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.WithField("address", cfg.Server.Http.Listen).Info("start http server")
		err := startHttpServer(ctx, cfg.Server.Http, router)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("server stopped")
				return
			}
			log.WithError(err).Error("error on http server")
		}
	}()

	wg.Wait()
	log.Info("application finished")
}

func startHttpServer(ctx context.Context, cfg *config.ServerConfig, router *chi.Mux) error {
	server := &http.Server{
		Addr:    cfg.Listen,
		Handler: router,
	}

	// cancel() всегда будет вызван, поэтому не используем waitgroup
	go func() {
		<-ctx.Done()
		err := server.Shutdown(context.Background())
		if err != nil {
			logger.New().WithError(err).Error("can't shutdown server")
		}
	}()

	return server.ListenAndServe()
}
