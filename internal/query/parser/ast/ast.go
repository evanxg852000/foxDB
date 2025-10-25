package ast

import (
	"fmt"
	"strings"

	"github.com/evanxg852000/foxdb/internal/types"
)

type Expression interface {
	ToExprString() string
}

type IdentifierExpr struct {
	Value string
}

func (ie *IdentifierExpr) ToExprString() string {
	return ie.Value
}

type StringLiteralExpr struct {
	Value string
}

func (sle *StringLiteralExpr) ToExprString() string {
	return "\"" + sle.Value + "\""
}

type IntegerLiteralExpr struct {
	Value int64
}

func (ile *IntegerLiteralExpr) ToExprString() string {
	return fmt.Sprintf("%d", ile.Value)
}

type FloatLiteralExpr struct {
	Value float64
}

func (fle *FloatLiteralExpr) ToExprString() string {
	return fmt.Sprintf("%f", fle.Value)
}

type NullLiteralExpr struct {
}

func (nle *NullLiteralExpr) ToExprString() string {
	return "NULL"
}

type BooleanLiteralExpr struct {
	Value bool
}

func (ble *BooleanLiteralExpr) ToExprString() string {
	return fmt.Sprintf("%t", ble.Value)
}

type AliasExpr struct {
	Alias string
	Expr  Expression
}

type CastExpr struct {
	Expr     Expression
	DataType string
}

type SortExpr struct {
	Expr      Expression
	Ascending bool
}

type PrefixExpr struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpr) ToExprString() string {
	return "(" + pe.Operator + pe.Right.ToExprString() + ")"
}

type InfixExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (be *InfixExpr) ToExprString() string {
	return "(" + be.Left.ToExprString() + " " + be.Operator + " " + be.Right.ToExprString() + ")"
}

type CallExpr struct {
	Function Expression
	Args     []Expression
}

func (ce *CallExpr) ToExprString() string {
	args := []string{}
	for _, arg := range ce.Args {
		args = append(args, arg.ToExprString())
	}
	return ce.Function.ToExprString() + "(" + strings.Join(args, ", ") + ")"
}

type Statement interface {
	ToStmtString() string
}

type Program struct {
	Statements []Statement
}

type CreateSchemaStatement struct {
	SchemaName  string
	IfNotExists bool
}

func (f *CreateSchemaStatement) ToStmtString() string {
	stmt := "CREATE SCHEMA "
	if f.IfNotExists {
		stmt += "IF NOT EXISTS "
	}
	stmt += f.SchemaName + ";"
	return stmt
}

type DropSchemaStatement struct {
	SchemaName string
}

func (d *DropSchemaStatement) ToStmtString() string {
	return "DROP SCHEMA " + d.SchemaName + ";"
}

type Constraint struct {
	PrimaryKey bool
	Unique     bool
	NotNull    bool
}
type ColumnDef struct {
	Name       string
	DataType   types.DataType // int, float, bool, text
	Constraint Constraint
}

type CreateTableStatement struct {
	TableName   string
	Columns     []ColumnDef
	IfNotExists bool
	PrimaryKeys []string
}

func (cts *CreateTableStatement) ToStmtString() string {
	stmt := "CREATE TABLE " + cts.TableName + " ("
	for i, col := range cts.Columns {
		stmt += col.Name + " " + col.DataType.String()
		if col.Constraint.PrimaryKey {
			stmt += " PRIMARY KEY"
		}
		if col.Constraint.Unique {
			stmt += " UNIQUE"
		}
		if col.Constraint.NotNull {
			stmt += " NOT NULL"
		}

		if i < len(cts.Columns)-1 {
			stmt += ", "
		}
	}

	if len(cts.PrimaryKeys) > 0 {
		stmt += " PRIMARY KEY (" + strings.Join(cts.PrimaryKeys, ", ") + ")"
	}

	stmt += ");"
	return stmt
}

type DropTableStatement struct {
	TableName string
}

func (dts *DropTableStatement) ToStmtString() string {
	return "DROP TABLE " + dts.TableName + ";"
}

type InsertStatement struct {
	TableName string
	Columns   []string
	Values    [][]Expression
}

type SelectStatement struct {
	Columns     []string
	FromClause  string
	WhereClause Expression
	GroupBy     []Expression
	OrderBy     []SortExpr
	Limit       uint64
	Offset      uint64
}
