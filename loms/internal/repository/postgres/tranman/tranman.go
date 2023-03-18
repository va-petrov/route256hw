package tranman

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type QueryEngine interface {
	Close()
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine
}

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

type txkey string

const key = txkey("tx")

func (tm *TransactionManager) RunTransaction(ctx context.Context, isoLevel pgx.TxIsoLevel, fx func(ctxTX context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx,
		pgx.TxOptions{
			IsoLevel: isoLevel,
		})
	if err != nil {
		return err
	}

	if err := fx(context.WithValue(ctx, key, tx)); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}

func (tm *TransactionManager) RunSerializable(ctx context.Context, fx func(ctxTX context.Context) error) error {
	return tm.RunTransaction(ctx, pgx.Serializable, fx)
}

func (tm *TransactionManager) RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error {
	return tm.RunTransaction(ctx, pgx.RepeatableRead, fx)
}

func (tm *TransactionManager) RunReadCommitted(ctx context.Context, fx func(ctxTX context.Context) error) error {
	return tm.RunTransaction(ctx, pgx.ReadCommitted, fx)
}

func (tm *TransactionManager) RunReadUncommitted(ctx context.Context, fx func(ctxTX context.Context) error) error {
	return tm.RunTransaction(ctx, pgx.ReadUncommitted, fx)
}

func (tm *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}
	return tm.pool
}
