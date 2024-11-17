package main

import (
	"log/slog"
	"referral-rest-api/internal/config"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/server"
	"referral-rest-api/internal/service"
	"referral-rest-api/internal/stopsignal"
	"referral-rest-api/internal/storage"
)

func main() {

	// Инициализируем конфиг и логгер.
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Debug("Config file and logger initialized", slog.String("env", cfg.Env))

	// Инициализируем пул подключений БД.
	db := storage.New(cfg)
	log.Debug("Storage initialized")

	// Инициализируем сервисный слой.
	app := service.New(cfg, log, db)
	log.Debug("Service initialized")

	// Инициализируем и запускаем HTTP сервер.
	srv := server.New(cfg)
	srv.Middleware()
	srv.API(app)
	srv.Start()
	log.Info("Server started")

	// Блокируем выполнение основной горутины до сигнала прерывания.
	stopsignal.Stop()

	// После сигнала прерывания останавливаем сервер.
	srv.GracefulShutdown()
	db.Close()
	log.Info("Server stopped")
}
