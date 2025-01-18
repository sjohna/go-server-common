package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sjohna/go-server-common/errors"
	"github.com/sjohna/go-server-common/log"
)

type Repo struct {
	DB *sqlx.DB
}

func (repo *Repo) NonTx(ctx context.Context) *DBDAO {
	return NewDBDAO(repo.DB, ctx)
}

func (repo *Repo) SerializableTx(ctx context.Context, transactionFunc func(*TxDAO) errors.Error) errors.Error {
	dao, err := NewTXDAO(repo.DB, ctx)
	if err != nil {

		return err
	}

	_, err = dao.Exec("set transaction isolation level serializable")
	if err != nil {
		if rollbackErr := dao.sqlxer.Rollback(); rollbackErr != nil {
			log.Ctx(ctx).WithError(errors.WrapDBError(rollbackErr, "failed to rollback transaction")).Error("Failed to rollback transaction!!!!")
		}
		return err
	}

	err = transactionFunc(dao)
	if err == nil {
		commitErr := dao.sqlxer.Commit()
		if commitErr != nil {
			wrappedCommitErr := errors.WrapDBError(commitErr, "failed to commit transaction")

			rollbackError := dao.sqlxer.Rollback()
			if rollbackError != nil {
				log.Ctx(dao.ctx).WithError(errors.WrapDBError(rollbackError, "failed to rollback transaction")).Error("Failed to rollback transaction after failing to commit!!!!")
			}

			return wrappedCommitErr
		}
	} else {
		rollbackError := dao.sqlxer.Rollback()
		if rollbackError != nil {
			log.Ctx(dao.ctx).WithError(errors.WrapDBError(rollbackError, "failed to rollback transaction")).Error("Failed to rollback transaction that returned an error!!!!")
		}
	}

	return err
}
