package physical

import (
	"context"

	"github.com/evanxg852000/foxdb/internal/query/planner"
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

func (p *UtilityPlan) Execute(ctx context.Context) (*types.DataChunk, error) {
	//TODO: implement execution logic
	return nil, nil
}
