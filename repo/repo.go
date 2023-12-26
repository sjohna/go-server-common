package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	DB *sqlx.DB
}

func (repo *Repo) NonTx(ctx context.Context) *DBDAO {
	return NewDBDAO(repo.DB, ctx)
}

func (repo *Repo) SerializableTx(ctx context.Context, transactionFunc func(*TxDAO) error) error {
	dao, err := NewTXDAO(repo.DB, ctx)
	logger := dao.Logger()
	if err != nil {
		logger.WithError(err).Error("Error creating TxDAO in SerializableTx")
		return err
	}

	_, err = dao.Exec("set transaction isolation level serializable")
	if err != nil {
		dao.Logger().WithError(err).Error("Failed to set transaction isolation level serializable")
		if rollbackErr := dao.sqlxer.Rollback(); err != nil {
			dao.Logger().WithField("seriousError", true).WithError(rollbackErr).Error("Failed to roll back transaction after failing to set isolation level serializable!!!")
		}
		return err
	}

	err = transactionFunc(dao)
	if err == nil {
		if err = dao.sqlxer.Commit(); err != nil {
			dao.Logger().WithError(err).Error("Failed to commit transaction!")

			if rollbackError := dao.sqlxer.Rollback(); err != nil {
				dao.Logger().WithField("seriousError", true).WithError(rollbackError).Error("Failled to roll back transaction after failing to commit!!!")
			}

			return err
		}
	} else {
		if err := dao.sqlxer.Rollback(); err != nil {
			dao.Logger().WithError(err).Error("Failled to roll back transaction.")
			return err
		}
	}

	return nil
}
