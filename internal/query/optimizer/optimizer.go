package optimizer

import (
	"context"

	"github.com/evanxg852000/foxdb/internal/catalog"
	"github.com/evanxg852000/foxdb/internal/query/optimizer/physical"
	"github.com/evanxg852000/foxdb/internal/query/planner"
	"github.com/evanxg852000/foxdb/internal/query/planner/logical"
	"github.com/evanxg852000/foxdb/internal/storage"
	"github.com/evanxg852000/foxdb/internal/types"
)

type PhysicalPlan interface {
	Execute(ctx context.Context, catalog *catalog.RootCatalog, storage *storage.KvStorage) (*types.DataChunk, error)
	GetSchema() *types.DataSchema
}

type Optimizer struct {
	catalog *catalog.RootCatalog
	stats   map[string]interface{}
}

func NewOptimizer(catalog *catalog.RootCatalog, stats map[string]interface{}) *Optimizer {
	return &Optimizer{
		catalog: catalog,
		stats:   stats,
	}
}

func (o *Optimizer) Optimize(logicalPlan planner.LogicalPlan) (PhysicalPlan, error) {
	//handle utility statements
	switch plan := logicalPlan.(type) {
	case *logical.CreateSchemaPlan:
		return physical.NewUtilityPlan(plan), nil
	}

	//TODO: implement a full optimization process
	schema := o.catalog.GetSchema("todo")
	table := schema.GetTable("todo")
	return physical.NewScan(table), nil
}
