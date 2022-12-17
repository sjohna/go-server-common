package repo

import (
	"context"
	"database/sql"
	"sync/atomic"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type DAO interface {
	Logger() *logrus.Entry

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
	logger *logrus.Entry
}

type TxDAO struct {
	sqlxer *sqlx.Tx
	logger *logrus.Entry
}

var daoIdCounter int64 = 0

func getNextDaoId() int64 {
	return atomic.AddInt64(&daoIdCounter, 1)
}

func NewDBDAO(db *sqlx.DB, logger *logrus.Entry) *DBDAO {
	if db == nil {
		logger.Fatal("db parameter not provided to NewDBDAO!")
	}

	if logger == nil {
		logger.Fatal("logger parameter not provided to NewDBDAO!")
	}

	log := logger.WithFields(logrus.Fields{
		"repo-dao-id": getNextDaoId(),
	})

	log.WithField("repo-dao-type", "non-tx").Info("DAO created")

	return &DBDAO{
		db,
		log,
	}
}

func (dao *DBDAO) Logger() *logrus.Entry {
	return dao.logger
}

func (dao *DBDAO) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dao.sqlxer.Exec(query, args...)
}

func (dao *DBDAO) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return dao.sqlxer.ExecContext(ctx, query, args...)
}

func (dao *DBDAO) Get(dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.Get(dest, query, args...)
}

func (dao *DBDAO) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.GetContext(ctx, dest, query, args...)
}

func (dao *DBDAO) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return dao.sqlxer.NamedExec(query, arg)
}

func (dao *DBDAO) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return dao.NamedExecContext(ctx, query, arg)
}

func (dao *DBDAO) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return dao.sqlxer.NamedQuery(query, arg)
}

func (dao *DBDAO) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return dao.sqlxer.PrepareNamed(query)
}

func (dao *DBDAO) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return dao.sqlxer.PrepareNamedContext(ctx, query)
}

func (dao *DBDAO) Preparex(query string) (*sqlx.Stmt, error) {
	return dao.sqlxer.Preparex(query)
}

func (dao *DBDAO) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return dao.sqlxer.PreparexContext(ctx, query)
}

func (dao *DBDAO) Rebind(query string) string {
	return dao.Rebind(query)
}

func (dao *DBDAO) Select(dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.Select(dest, query, args...)
}

func (dao *DBDAO) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.SelectContext(ctx, dest, query, args...)
}

func (dao *DBDAO) Unsafe() DAO {
	return &DBDAO{
		dao.sqlxer.Unsafe(),
		dao.logger.WithField("repo-unsafe", true),
	}
}

func NewTXDAO(db *sqlx.DB, logger *logrus.Entry) (*TxDAO, error) {
	if db == nil {
		logger.Fatal("db parameter not provided to NewTXDAO!")
	}

	if logger == nil {
		logger.Fatal("logger parameter not provided to NewTXDAO!")
	}

	log := logger.WithFields(logrus.Fields{
		"repo-dao-id": getNextDaoId(),
	})

	tx, err := db.Beginx()
	if err != nil {
		log.WithField("repo-dao-type", "tx").WithError(err).Error("Error beginning transaction")
		return nil, err
	}

	log.WithField("repo-dao-type", "tx").Info("DAO created")

	return &TxDAO{
		tx,
		log,
	}, nil
}

func (dao *TxDAO) Logger() *logrus.Entry {
	return dao.logger
}

func (dao *TxDAO) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dao.sqlxer.Exec(query, args...)
}

func (dao *TxDAO) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return dao.sqlxer.ExecContext(ctx, query, args...)
}

func (dao *TxDAO) Get(dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.Get(dest, query, args...)
}

func (dao *TxDAO) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.GetContext(ctx, dest, query, args...)
}

func (dao *TxDAO) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return dao.sqlxer.NamedExec(query, arg)
}

func (dao *TxDAO) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return dao.NamedExecContext(ctx, query, arg)
}

func (dao *TxDAO) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return dao.sqlxer.NamedQuery(query, arg)
}

func (dao *TxDAO) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return dao.sqlxer.PrepareNamed(query)
}

func (dao *TxDAO) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return dao.sqlxer.PrepareNamedContext(ctx, query)
}

func (dao *TxDAO) Preparex(query string) (*sqlx.Stmt, error) {
	return dao.sqlxer.Preparex(query)
}

func (dao *TxDAO) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return dao.sqlxer.PreparexContext(ctx, query)
}

func (dao *TxDAO) Rebind(query string) string {
	return dao.Rebind(query)
}

func (dao *TxDAO) Select(dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.Select(dest, query, args...)
}

func (dao *TxDAO) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dao.sqlxer.SelectContext(ctx, dest, query, args...)
}

func (dao *TxDAO) Unsafe() DAO {
	return &TxDAO{
		dao.sqlxer.Unsafe(),
		dao.logger.WithField("repo-unsafe", true),
	}
}
