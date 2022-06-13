package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"

	logger "github.com/sirupsen/logrus"

	"github.com/xRepZ/payments/cmd/config"

	"github.com/go-chi/chi/v5"
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

	// db, err := sqlx.Connect("postgres", cfg.Storage.Postgres.Dsn)
	// if err != nil {
	// 	log.WithError(err).Fatal("can't connect to postgres")
	// }

	// Объявляем handlers
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome anonymous"))
		})
		// Вызывается пользователем
		// r.Post("/api/transaction",) // Создание платежа
		// r.Get("/api/transaction/{transaction_id}",) // Проверка статуса платежа по id
		// r.Get("/api/transaction/user/{user_id}", ) // Все платежи юзера по его id
		// r.Get("/api/transaction/user/{user_email}", ) // Все платежи юзера по его email
		// r.Delete("/api/transaction/{transaction_id}", ) // Отмена платежа по его id
		// Вызывается платежной системой
		// r.With(AuthMiddleware).Patch("/api/transaction/{transaction_id}", ) // Изменение статуса платежа
	})

	wg := sync.WaitGroup{}

	wg.Add(1)
	osSigCh := make(chan os.Signal, 1)
	signal.Notify(osSigCh, os.Interrupt, os.Kill)
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
