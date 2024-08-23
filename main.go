package main

import (
	"log/slog"
	"os"

	"github.com/Foxtrot1388/MessaggioTask/internal/config"
	"github.com/Foxtrot1388/MessaggioTask/internal/server"
	"github.com/Foxtrot1388/MessaggioTask/internal/service"
	kafka "github.com/Foxtrot1388/MessaggioTask/internal/storage/kafka"
	storagepg "github.com/Foxtrot1388/MessaggioTask/internal/storage/pg"
)

// @title Message API
// @version 1.0
// @description API Server for messages

// @host 89.169.167.99:8080
// @BasePath /

func main() {

	cfg := config.GetConfig()
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(cfg.LogLevel)}),
	)
	log.Info("Use config", "data", cfg)

	db, err := storagepg.New(cfg, log)
	if err != nil {
		panic(err)
	}

	kafkastorage, err := kafka.New(cfg, log)
	if err != nil {
		panic(err)
	}
	usercases := service.New(log, db, kafkastorage)
	app := server.New(log, usercases)

	app.Listen()

}
