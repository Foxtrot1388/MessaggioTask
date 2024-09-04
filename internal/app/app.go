package app

import (
	"log/slog"
	"os"

	"github.com/Foxtrot1388/MessaggioTask/internal/config"
	server "github.com/Foxtrot1388/MessaggioTask/internal/controller/http"
	"github.com/Foxtrot1388/MessaggioTask/internal/service"
	kafka "github.com/Foxtrot1388/MessaggioTask/internal/storage/kafka"
	storagepg "github.com/Foxtrot1388/MessaggioTask/internal/storage/pg"
)

func Run() {

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
