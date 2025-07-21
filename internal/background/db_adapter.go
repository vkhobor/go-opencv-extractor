package background

import (
	"context"
	"database/sql"

	"github.com/vkhobor/go-opencv/internal/db"
	qDB "github.com/vkhobor/go-opencv/internal/db"
	"github.com/vkhobor/go-opencv/internal/features"
)

type DdAdapter struct {
	Querier features.QuerierWithTx
	TxEr    features.TXer
}

func NewDbAdapter(db *sql.DB) DdAdapter {
	return DdAdapter{
		Querier: newQuerierAdapter(db),
		TxEr:    newTxAdapter(db),
	}
}

type querierAdapter struct {
	db.Queries
}

// WithTx implements features.QuerierWithTx.
func (q *querierAdapter) WithTx(tx db.DBTX) features.QuerierWithTx {
	return &querierAdapter{Queries: *qDB.New(tx)}
}

var _ features.QuerierWithTx = &querierAdapter{}

func newQuerierAdapter(db *sql.DB) *querierAdapter {
	return &querierAdapter{Queries: *qDB.New(db)}
}

type txAdapter struct {
	db *sql.DB
}

// BeginTx implements features.TXer.
func (t txAdapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (features.TX, error) {
	return t.db.BeginTx(ctx, opts)
}

var _ features.TXer = txAdapter{}

func newTxAdapter(db *sql.DB) *txAdapter {
	return &txAdapter{db: db}
}
