package uow

import "database/sql"

type IUnitOfWork interface {
	Do(fn func(tx *sql.Tx) error) error
}
