package physical

import (
	"context"

	"github.com/evanxg852000/foxdb/internal/catalog"
	"github.com/evanxg852000/foxdb/internal/query/planner"
	"github.com/evanxg852000/foxdb/internal/query/planner/logical"
	"github.com/evanxg852000/foxdb/internal/storage"
	"github.com/evanxg852000/foxdb/internal/types"
)

// A generic plan for executing utility statements
// This is just like a wrapper around logical plan that
// should be executed as is without any optimization
type UtilityPlan struct {
	logicalPlan planner.LogicalPlan
}

func NewUtilityPlan(logicalPlan planner.LogicalPlan) *UtilityPlan {
	return &UtilityPlan{
		logicalPlan: logicalPlan,
	}
}

func (p *UtilityPlan) GetSchema() *types.DataSchema {
	return p.logicalPlan.GetSchema()
}

func (p *UtilityPlan) Execute(ctx context.Context, catalog *catalog.RootCatalog, storage *storage.KvStorage) (*types.DataChunk, error) {
	//TODO: complete execution logic for utility plans
	switch plan := p.logicalPlan.(type) {
	case *logical.CreateSchemaPlan:
		return createSchema(catalog, plan.SchemaName, plan.IfNotExists)
	}
	return nil, nil
}

func createSchema(catalog *catalog.RootCatalog, name string, safe bool) (*types.DataChunk, error) {
	catalog.Lock()
	defer catalog.Unlock()
	if safe {
		if schema := catalog.GetSchema(name); schema != nil {
			return nil, nil
		}
	}

	if _, err := catalog.AddSchema(name); err != nil {
		return nil, err
	}
	return nil, nil
}
