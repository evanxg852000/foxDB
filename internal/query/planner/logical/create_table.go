package logical

import (
	"github.com/evanxg852000/foxdb/internal/catalog"
	"github.com/evanxg852000/foxdb/internal/query/parser/ast"
	"github.com/evanxg852000/foxdb/internal/types"
)

type CreateTablePlan struct {
	TableName   string
	Columns     []catalog.Column
	IfNotExists bool
	PrimaryKeys []string
}

func NewCreateTablePlan(statement *ast.CreateTableStatement) *CreateTablePlan {
	columns := make([]catalog.Column, 0, len(statement.Columns))
	for _, colDef := range statement.Columns {
		constraints := catalog.Constraint{
			Unique:  colDef.Constraint.Unique,
			NotNull: colDef.Constraint.NotNull,
		}
		column := catalog.NewColumn(0, colDef.Name, colDef.DataType, constraints)
		columns = append(columns, *column)
	}

	return &CreateTablePlan{
		TableName:   statement.TableName,
		Columns:     columns,
		IfNotExists: statement.IfNotExists,
		PrimaryKeys: statement.PrimaryKeys,
	}
}

func (p *CreateTablePlan) GetSchema() *types.DataSchema {
	return nil
}
