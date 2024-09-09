package service

import (
	"context"
	"time"

	"github.com/Foxtrot1388/MessaggioTask/internal/entity"
	"github.com/Foxtrot1388/MessaggioTask/internal/model"
)

//go:generate mockery --name DbRepository
type DbRepository interface {
	Create(context.Context, []string) ([]entity.OutputMessage, error)
	CreateOutbox(context.Context, []model.OutputMessage) error
	SelectForOutbox(context.Context) ([]entity.OutputMessageOutbox, error)
	DeleteForOutboxByID(context.Context, int) error
	UpdateForOutboxByID(context.Context, int) error
	AddStatistic(context.Context, int) error
	SelectStatistic(context.Context, time.Time, time.Time) ([]entity.StatMessage, error)
	WithTr(context.Context, func(context.Context) (interface{}, error)) (interface{}, error)
}

//go:generate mockery --name KafkaRepository
type KafkaRepository interface {
	Send(int) error
	Read() (<-chan int, error)
}
