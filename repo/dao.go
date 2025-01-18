package repo

import (
	"context"
	"database/sql"
	"github.com/sjohna/go-server-common/errors"
	"github.com/sjohna/go-server-common/log"
	"sync/atomic"

	"github.com/jmoiron/sqlx"
)

type DAO interface {
	Context() context.Context
	Exec(query string, args ...interface{}) (sql.Result, errors.Error)
	Get(dest interface{}, query string, args ...interface{}) errors.Error
	NamedExec(query string, arg interface{}) (sql.Result, errors.Error)
	PrepareNamed(query string) (*sqlx.NamedStmt, errors.Error)
	Preparex(query string) (*sqlx.Stmt, errors.Error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) errors.Error
	Unsafe() DAO
}

// might be able to refactor this so I don't need to duplicate all the functions
type DBDAO struct {
	sqlxer *sqlx.DB
	ctx    context.Context
}

type TxDAO struct {
	sqlxer *sqlx.Tx
	ctx    context.Context
}

var daoIdCounter int64 = 0

func getNextDaoId() int64 {
	return atomic.AddInt64(&daoIdCounter, 1)
}

func NewDBDAO(db *sqlx.DB, ctx context.Context) *DBDAO {
	logger := ctx.Value("logger").(log.Logger)

	if logger == nil {
		panic("logger not in context provided to NewDBDAO!")
	}

	if db == nil {
		logger.Panic("db parameter not provided to NewDBDAO!")
		panic("db parameter not provided to NewDBDAO!")
	}

	DAOLogger := logger.WithField("repo-dao-id", getNextDaoId())

	DAOLogger.WithField("repo-dao-type", "non-tx").Debug("DAO created")
	DAOCtx := context.WithValue(ctx, "logger", DAOLogger)

	return &DBDAO{
		db,
		DAOCtx,
	}
}

func (dao *DBDAO) Context() context.Context {
	return dao.ctx
}

func (dao *DBDAO) Exec(query string, args ...interface{}) (sql.Result, errors.Error) {
	var myErr errors.Error
	result, err := dao.sqlxer.ExecContext(dao.ctx, query, args...)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running Exec", query, args)
	}
	return result, myErr
}

func (dao *DBDAO) Get(dest interface{}, query string, args ...interface{}) errors.Error {
	var myErr errors.Error
	err := dao.sqlxer.GetContext(dao.ctx, dest, query, args...)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running Get", query, args)
	}
	return myErr
}

func (dao *DBDAO) NamedExec(query string, arg interface{}) (sql.Result, errors.Error) {
	var myErr errors.Error
	result, err := dao.sqlxer.NamedExecContext(dao.ctx, query, arg)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running NamedExec", query, arg)
	}
	return result, myErr
}

func (dao *DBDAO) PrepareNamed(query string) (*sqlx.NamedStmt, errors.Error) {
	var myErr errors.Error
	namedStmnt, err := dao.sqlxer.PrepareNamedContext(dao.ctx, query)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running PrepareNamed", query, nil)
	}
	return namedStmnt, myErr
}

func (dao *DBDAO) Preparex(query string) (*sqlx.Stmt, errors.Error) {
	var myErr errors.Error
	stmnt, err := dao.sqlxer.PreparexContext(dao.ctx, query)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running Preparex", query, nil)
	}
	return stmnt, myErr
}

func (dao *DBDAO) Rebind(query string) string {
	return dao.sqlxer.Rebind(query)
}

func (dao *DBDAO) Select(dest interface{}, query string, args ...interface{}) errors.Error {
	var myErr errors.Error
	err := dao.sqlxer.SelectContext(dao.ctx, dest, query, args...)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running Select", query, args)
	}
	return myErr
}

func (dao *DBDAO) Unsafe() DAO {
	logger := log.Ctx(dao.ctx).WithField("repo-unsafe", true)
	newCtx := context.WithValue(dao.ctx, "logger", logger)
	logger.Info("Unsafe DBDAO created")
	return &DBDAO{
		dao.sqlxer.Unsafe(),
		newCtx,
	}
}

func NewTXDAO(db *sqlx.DB, ctx context.Context) (*TxDAO, errors.Error) {
	if db == nil {
		log.Ctx(ctx).Panic("db parameter not provided to NewTXDAO!")
		panic("db parameter not provided to NewTXDAO!")
	}

	txLogger := log.Ctx(ctx).WithField("repo-dao-id", getNextDaoId())

	tx, err := db.Beginx()
	if err != nil {
		myErr := errors.Wrap(err, "Error beginning transaction")
		return nil, myErr
	}

	txLogger.WithField("repo-dao-type", "tx").Debug("TXDAO created")
	txCtx := context.WithValue(ctx, "logger", txLogger)

	return &TxDAO{
		tx,
		txCtx,
	}, nil
}

func (dao *TxDAO) Context() context.Context {
	return dao.ctx
}

func (dao *TxDAO) Exec(query string, args ...interface{}) (sql.Result, errors.Error) {
	var myErr errors.Error
	result, err := dao.sqlxer.ExecContext(dao.ctx, query, args...)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running Exec", query, args)
	}
	return result, myErr
}

func (dao *TxDAO) Get(dest interface{}, query string, args ...interface{}) errors.Error {
	var myErr errors.Error
	err := dao.sqlxer.GetContext(dao.ctx, dest, query, args...)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running Get", query, args)
	}
	return myErr
}

func (dao *TxDAO) NamedExec(query string, arg interface{}) (sql.Result, errors.Error) {
	var myErr errors.Error
	result, err := dao.sqlxer.NamedExecContext(dao.ctx, query, arg)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running NamedExec", query, arg)
	}
	return result, myErr
}

func (dao *TxDAO) PrepareNamed(query string) (*sqlx.NamedStmt, errors.Error) {
	var myErr errors.Error
	namedStmnt, err := dao.sqlxer.PrepareNamedContext(dao.ctx, query)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running PrepareNamed", query, nil)
	}
	return namedStmnt, myErr
}

func (dao *TxDAO) Preparex(query string) (*sqlx.Stmt, errors.Error) {
	var myErr errors.Error
	stmnt, err := dao.sqlxer.PreparexContext(dao.ctx, query)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running Preparex", query, nil)
	}
	return stmnt, myErr
}

func (dao *TxDAO) Rebind(query string) string {
	return dao.sqlxer.Rebind(query)
}

func (dao *TxDAO) Select(dest interface{}, query string, args ...interface{}) errors.Error {
	var myErr errors.Error
	err := dao.sqlxer.SelectContext(dao.ctx, dest, query, args...)
	if err != nil {
		myErr = errors.WrapQueryError(err, "Error running Select", query, args)
	}
	return myErr
}

func (dao *TxDAO) Unsafe() DAO {
	logger := log.Ctx(dao.ctx).WithField("repo-unsafe", true)
	newCtx := context.WithValue(dao.ctx, "logger", logger)
	logger.Info("Unsafe TxDAO created")
	return &TxDAO{
		dao.sqlxer.Unsafe(),
		newCtx,
	}
}
