package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Foxtrot1388/MessaggioTask/internal/model"
	"github.com/Foxtrot1388/MessaggioTask/internal/model/converter"
)

type service struct {
	db    DbRepository
	kafka KafkaRepository
	log   *slog.Logger
}

func New(log *slog.Logger, db DbRepository, kafka KafkaRepository) *service {
	return &service{db: db, log: log, kafka: kafka}
}

func (s *service) Create(ctx context.Context, messages []string) ([]model.OutputMessage, error) {

	// в транзакции запишем в аутбокс и в сообщения
	res, err := s.db.WithTr(ctx, func(ctx context.Context) (interface{}, error) {

		entityresults, err := s.db.Create(ctx, messages)
		if err != nil {
			return nil, err
		}

		result := make([]model.OutputMessage, len(entityresults))
		for c, entityresult := range entityresults {
			result[c] = converter.GetOutputMessage(&entityresult)
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
		return res.([]model.OutputMessage), nil
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

func (s *service) GetStatistic(ctx context.Context, dateAt, dateTo time.Time) ([]model.StatMessage, error) {

	entityresult, err := s.db.SelectStatistic(ctx, dateAt, dateTo)

	result := make([]model.StatMessage, len(entityresult))
	for c, element := range entityresult {
		result[c] = converter.GetStatMessage(&element)
	}

	return result, err

}
