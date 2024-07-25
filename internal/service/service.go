package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Foxtrot1388/MessaggioTask/internal/entity"
)

//go:generate mockery --name DbRepository
type DbRepository interface {
	Create(context.Context, []entity.MessageToInsert) ([]entity.OutputMessage, error)
	CreateOutbox(context.Context, []entity.OutputMessage) error
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

type service struct {
	db    DbRepository
	kafka KafkaRepository
	log   *slog.Logger
}

func New(log *slog.Logger, db DbRepository, kafka KafkaRepository) *service {
	return &service{db: db, log: log, kafka: kafka}
}

func (s *service) Create(ctx context.Context, messages []entity.MessageToInsert) ([]entity.OutputMessage, error) {

	// в транзакции запишем в аутбокс и в сообщения
	res, err := s.db.WithTr(ctx, func(ctx context.Context) (interface{}, error) {

		result, err := s.db.Create(ctx, messages)
		if err != nil {
			return nil, err
		}

		err = s.db.CreateOutbox(ctx, result)
		if err != nil {
			return nil, err
		}

		return result, nil

	})
	if err != nil {
		return nil, err
	} else {
		return res.([]entity.OutputMessage), nil
	}

}

func (s *service) StartJobOutboxWrite(ctx context.Context) {

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		default:
			_ = <-ticker.C
			res, err := s.db.SelectForOutbox(ctx)
			if err != nil {
				s.log.Error(err.Error())
				continue
			}
			for _, v := range res {
				err := s.kafka.Send(v.Idmessage)
				if err == nil {
					err = s.db.DeleteForOutboxByID(ctx, v.ID)
					if err != nil {
						s.log.Error(err.Error())
						continue
					}
				} else {
					s.log.Error(err.Error())
					continue
				}
			}
		}
	}

}

func (s *service) StartJobOutboxRead(ctx context.Context) {

	ch, err := s.kafka.Read()
	if err != nil {
		s.log.Error(err.Error())
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			err = s.db.UpdateForOutboxByID(ctx, msg)
			_ = s.db.AddStatistic(ctx, msg)
			if err != nil {
				s.log.Error(err.Error())
				continue
			}
		}
	}

}

func (s *service) GetStatistic(ctx context.Context, dateAt, dateTo time.Time) ([]entity.StatMessage, error) {

	return s.db.SelectStatistic(ctx, dateAt, dateTo)

}
