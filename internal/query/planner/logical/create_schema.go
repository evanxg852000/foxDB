package logical

import (
	"github.com/evanxg852000/foxdb/internal/query/parser/ast"
	"github.com/evanxg852000/foxdb/internal/types"
)

type CreateSchemaPlan struct {
	SchemaName  string
	IfNotExists bool
}

func NewCreateSchemaPlan(statement *ast.CreateSchemaStatement) *CreateSchemaPlan {
	return &CreateSchemaPlan{
		SchemaName:  statement.SchemaName,
		IfNotExists: statement.IfNotExists,
	}
}

func (p *CreateSchemaPlan) GetSchema() *types.DataSchema {
	return nil
}
