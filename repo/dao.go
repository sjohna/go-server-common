package repo

import (
	"context"
	"database/sql"
	"github.com/sjohna/go-server-common/log"
	"sync/atomic"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type DAO interface {
	Context() context.Context
	Logger() log.Logger

	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Unsafe() DAO
}

type DBDAO struct {
	sqlxer *sqlx.DB
	ctx    context.Context
	logger log.Logger
}

type TxDAO struct {
	sqlxer *sqlx.Tx
	ctx    context.Context
	logger log.Logger
}

var daoIdCounter int64 = 0

func getNextDaoId() int64 {
	return atomic.AddInt64(&daoIdCounter, 1)
}

func NewDBDAO(db *sqlx.DB, ctx context.Context) *DBDAO {
	logger := ctx.Value("logger").(log.Logger)
	if db == nil {
		logger.Fatal("db parameter not provided to NewDBDAO!")
	}

	if logger == nil {
		logger.Fatal("logger parameter not provided to NewDBDAO!")
	}

	DAOLogger := logger.WithFields(logrus.Fields{
		"repo-dao-id": getNextDaoId(),
	})

	DAOLogger.WithField("repo-dao-type", "non-tx").Info("DAO created")
	DAOCtx := context.WithValue(ctx, "logger", DAOLogger)

	return &DBDAO{
		db,
		DAOCtx,
		DAOLogger,
	}
}

func (dao *DBDAO) Logger() log.Logger {
	return dao.logger
}

func (dao *DBDAO) Context() context.Context {
	return dao.ctx
}

func (dao *DBDAO) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dao.sqlxer.ExecContext(dao.ctx, query, args...)
}

func (dao *DBDAO) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return dao.sqlxer.ExecContext(ctx, query, args...)
}

func (dao *DBDAO) Get(dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.GetContext(dao.ctx, dest, query, args...)
}

func (dao *DBDAO) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.GetContext(ctx, dest, query, args...)
}

func (dao *DBDAO) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return dao.sqlxer.NamedExecContext(dao.ctx, query, arg)
}

func (dao *DBDAO) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return dao.NamedExecContext(ctx, query, arg)
}

func (dao *DBDAO) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return dao.sqlxer.NamedQueryContext(dao.ctx, query, arg)
}

func (dao *DBDAO) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return dao.sqlxer.PrepareNamedContext(dao.ctx, query)
}

func (dao *DBDAO) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return dao.sqlxer.PrepareNamedContext(ctx, query)
}

func (dao *DBDAO) Preparex(query string) (*sqlx.Stmt, error) {
	return dao.sqlxer.PreparexContext(dao.ctx, query)
}

func (dao *DBDAO) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return dao.sqlxer.PreparexContext(ctx, query)
}

func (dao *DBDAO) Rebind(query string) string {
	return dao.Rebind(query)
}

func (dao *DBDAO) Select(dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.SelectContext(dao.ctx, dest, query, args...)
}

func (dao *DBDAO) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.SelectContext(ctx, dest, query, args...)
}

func (dao *DBDAO) Unsafe() DAO {
	logger := dao.logger.WithField("repo-unsafe", true)
	newCtx := context.WithValue(dao.ctx, "logger", logger)
	logger.Info("Unsafe DBDAO created")
	return &DBDAO{
		dao.sqlxer.Unsafe(),
		newCtx,
		logger,
	}
}

func NewTXDAO(db *sqlx.DB, ctx context.Context) (*TxDAO, error) {
	logger := ctx.Value("logger").(log.Logger)

	if db == nil {
		logger.Fatal("db parameter not provided to NewTXDAO!")
	}

	if logger == nil {
		logger.Fatal("logger parameter not provided to NewTXDAO!")
	}

	txLogger := logger.WithFields(logrus.Fields{
		"repo-dao-id": getNextDaoId(),
	})

	tx, err := db.Beginx()
	if err != nil {
		txLogger.WithField("repo-dao-type", "tx").WithError(err).Error("Error beginning transaction")
		return nil, err
	}

	txLogger.WithField("repo-dao-type", "tx").Info("TXDAO created")
	txCtx := context.WithValue(ctx, "logger", txLogger)

	return &TxDAO{
		tx,
		txCtx,
		txLogger,
	}, nil
}

func (dao *TxDAO) Logger() log.Logger {
	return dao.logger
}

func (dao *TxDAO) Context() context.Context {
	return dao.ctx
}

func (dao *TxDAO) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dao.sqlxer.ExecContext(dao.ctx, query, args...)
}

func (dao *TxDAO) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return dao.sqlxer.ExecContext(ctx, query, args...)
}

func (dao *TxDAO) Get(dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.GetContext(dao.ctx, dest, query, args...)
}

func (dao *TxDAO) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.GetContext(ctx, dest, query, args...)
}

func (dao *TxDAO) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return dao.sqlxer.NamedExecContext(dao.ctx, query, arg)
}

func (dao *TxDAO) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return dao.NamedExecContext(ctx, query, arg)
}

func (dao *TxDAO) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return dao.sqlxer.NamedQuery(query, arg)
}

func (dao *TxDAO) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return dao.sqlxer.PrepareNamedContext(dao.ctx, query)
}

func (dao *TxDAO) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return dao.sqlxer.PrepareNamedContext(ctx, query)
}

func (dao *TxDAO) Preparex(query string) (*sqlx.Stmt, error) {
	return dao.sqlxer.PreparexContext(dao.ctx, query)
}

func (dao *TxDAO) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return dao.sqlxer.PreparexContext(ctx, query)
}

func (dao *TxDAO) Rebind(query string) string {
	return dao.Rebind(query)
}

func (dao *TxDAO) Select(dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.SelectContext(dao.ctx, dest, query, args...)
}

func (dao *TxDAO) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.SelectContext(ctx, dest, query, args...)
}

func (dao *TxDAO) Unsafe() DAO {
	logger := dao.logger.WithField("repo-unsafe", true)
	newCtx := context.WithValue(dao.ctx, "logger", logger)
	logger.Info("Unsafe TxDAO created")
	return &TxDAO{
		dao.sqlxer.Unsafe(),
		newCtx,
		logger,
	}
}
