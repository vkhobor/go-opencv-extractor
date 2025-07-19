package features

import (
	"context"
	"database/sql"

	"github.com/vkhobor/go-opencv/db"
)

type QuerierWithTx interface {
	db.Querier
	WithTx(tx db.DBTX) QuerierWithTx
}

type TXer interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (TX, error)
}

type TX interface {
	db.DBTX
	Commit() error
	Rollback() error
}
