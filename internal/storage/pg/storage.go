package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/Foxtrot1388/MessaggioTask/internal/config"
	"github.com/Foxtrot1388/MessaggioTask/internal/entity"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type dbStorage struct {
	log *slog.Logger
	cfg *config.AppConfig
	DB  *sqlx.DB
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func New(cfg *config.AppConfig, log *slog.Logger) (*dbStorage, error) {

	db, err := sqlx.Open("pgx", cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	result := &dbStorage{log: log, cfg: cfg, DB: db}

	err = result.migrateDB()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (con *dbStorage) WithTr(ctx context.Context, f func(context.Context) (interface{}, error)) (interface{}, error) {

	tx, err := con.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	ctxtx := context.WithValue(ctx, "tx", tx)

	result, err := f(ctxtx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return result, nil

}

func (con *dbStorage) Create(ctx context.Context, messages []entity.MessageToInsert) ([]entity.OutputMessage, error) {

	tx := ctx.Value("tx").(*sqlx.Tx)

	op := psql.
		Insert("messages").
		Columns("Message", "processed").
		Suffix("RETURNING \"id\"")
	for _, message := range messages {
		op = op.Values(message.Message, false)
	}

	query, args, err := op.ToSql()
	if err != nil {
		return nil, err
	}
	var queryer sqlx.QueryerContext
	if tx != nil {
		queryer = tx
	} else {
		queryer = con.DB
	}
	res := make([]entity.OutputMessage, 0)
	rows, err := queryer.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item entity.OutputMessage
		err = rows.Scan(
			&item.ID,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	return res, nil

}

func (con *dbStorage) CreateOutbox(ctx context.Context, messages []entity.OutputMessage) error {

	tx := ctx.Value("tx").(*sqlx.Tx)

	op := psql.
		Insert("outboxmessages").
		Columns("idmessage")
	for _, message := range messages {
		op = op.Values(message.ID)
	}

	query, args, err := op.ToSql()
	if err != nil {
		return err
	}
	var queryer sqlx.QueryerContext
	if tx != nil {
		queryer = tx
	} else {
		queryer = con.DB
	}
	_, err = queryer.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	} else {
		return nil
	}

}

func (con *dbStorage) SelectForOutbox(ctx context.Context) ([]entity.OutputMessageOutbox, error) {

	op := psql.
		Select("ID", "idmessage").
		From("outboxmessages")

	query, args, err := op.ToSql()
	if err != nil {
		return nil, err
	}

	res := make([]entity.OutputMessageOutbox, 0)
	rows, err := con.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item entity.OutputMessageOutbox
		err = rows.Scan(
			&item.ID,
			&item.Idmessage,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	return res, nil
}

func (con *dbStorage) DeleteForOutboxByID(ctx context.Context, id int) error {

	op := psql.
		Delete("").
		From("outboxmessages").
		Where("ID = ?", id)

	query, args, err := op.ToSql()
	if err != nil {
		return err
	}

	_, err = con.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil

}

func (con *dbStorage) UpdateForOutboxByID(ctx context.Context, id int) error {

	op := psql.
		Update("messages").
		Set("processed", true).
		Where("id = ?", id)

	query, args, err := op.ToSql()
	if err != nil {
		return err
	}

	_, err = con.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil

}

func (con *dbStorage) AddStatistic(ctx context.Context, id int) error {

	op := psql.
		Insert("agregate_processed_messages").
		Columns("date", "id").
		Values(time.Now(), id)

	query, args, err := op.ToSql()
	if err != nil {
		return err
	}

	_, err = con.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil

}

func (con *dbStorage) SelectStatistic(ctx context.Context, dateAt, dateTo time.Time) ([]entity.StatMessage, error) {

	op := psql.
		Select("count(id) id", "date_trunc('day', date) date").
		From("agregate_processed_messages").
		GroupBy("date_trunc('day', date)")

	query, args, err := op.ToSql()
	if err != nil {
		return nil, err
	}

	res := make([]entity.StatMessage, 0)
	rows, err := con.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item entity.StatMessage
		err = rows.Scan(
			&item.Count,
			&item.Day,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	return res, nil

}
