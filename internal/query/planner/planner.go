package planner

import (
	"fmt"

	"github.com/evanxg852000/foxdb/internal/catalog"
	"github.com/evanxg852000/foxdb/internal/query/parser/ast"
	"github.com/evanxg852000/foxdb/internal/query/planner/logical"
	"github.com/evanxg852000/foxdb/internal/types"
)

type LogicalPlan interface {
	GetSchema() *types.DataSchema
}

// plan and bind the ast to generate a logical plan
type Planner struct {
	catalog *catalog.RootCatalog
}

func NewPlanner(catalog *catalog.RootCatalog) *Planner {
	return &Planner{
		catalog: catalog,
	}
}

func (p *Planner) Plan(queryAst ast.Statement) (LogicalPlan, error) {
	switch stmt := queryAst.(type) {
	case *ast.CreateSchemaStatement:
		return logical.NewCreateSchemaPlan(stmt), nil

	default:
		return nil, fmt.Errorf("unsupported statement type: %T", queryAst)
	}
}
