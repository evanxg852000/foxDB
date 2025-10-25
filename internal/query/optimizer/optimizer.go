package optimizer

import (
	"context"

	"github.com/evanxg852000/foxdb/internal/catalog"
	"github.com/evanxg852000/foxdb/internal/query/optimizer/physical"
	"github.com/evanxg852000/foxdb/internal/query/planner"
	"github.com/evanxg852000/foxdb/internal/types"
)

type PhysicalPlan interface {
	Execute(ctx context.Context) (*types.DataChunk, error)
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
	//TODO: implement a full optimization process
	schema, _ := o.catalog.GetSchema("todo")
	table, _ := schema.GetTable("todo")
	return physical.NewScanExec(table), nil
}
