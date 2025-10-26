package executor

import (
	"context"

	"github.com/evanxg852000/foxdb/internal/catalog"
	"github.com/evanxg852000/foxdb/internal/query/optimizer"
	"github.com/evanxg852000/foxdb/internal/storage"
	"github.com/evanxg852000/foxdb/internal/types"
)

type Executor struct {
	storage *storage.KvStorage
	catalog *catalog.RootCatalog
	plan    optimizer.PhysicalPlan
}

func NewExecutor(storage *storage.KvStorage, catalog *catalog.RootCatalog, plan optimizer.PhysicalPlan) *Executor {
	return &Executor{
		storage: storage,
		catalog: catalog,
		plan:    plan,
	}
}

func (e *Executor) Execute(ctx context.Context) (*types.DataChunk, error) {
	return e.plan.Execute(ctx, e.catalog, e.storage)
}
