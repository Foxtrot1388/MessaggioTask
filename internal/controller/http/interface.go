package server

import (
	"context"
	"time"

	"github.com/Foxtrot1388/MessaggioTask/internal/model"
)

type usecases interface {
	Create(context.Context, []string) ([]model.OutputMessage, error)
	GetStatistic(context.Context, time.Time, time.Time) ([]model.StatMessage, error)
	StartJobOutboxWrite(context.Context)
	StartJobOutboxRead(context.Context)
}
