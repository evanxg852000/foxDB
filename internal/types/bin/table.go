package types

import (
	"fmt"
	"slices"

	"github.com/evanxg852000/foxdb/internal/common"
)

type Constraint struct {
	PrimaryKey bool
	Unique     bool
	NotNull    bool
}

type Column struct {
	Id          common.Id
	Name        string
	DataType    DataType
	Constraints Constraint
}

type Index struct {
	Id        common.Id
	Name      string
	ColumnIds []common.Id
	Unique    bool
}

type TableDesc struct {
	Columns     []*Column
	ColumnNames map[string]int
	Indexes     []*Index
	IndexNames  map[string]int
}

func NewTableDesc() *TableDesc {
	return &TableDesc{
		Columns:     make([]*Column, 0),
		ColumnNames: make(map[string]int),
		Indexes:     make([]*Index, 0),
		IndexNames:  make(map[string]int),
	}
}

func (t *TableDesc) AddColumn(column *Column) error {
	_, exists := t.ColumnNames[column.Name]
	if exists {
		return fmt.Errorf("column %s already exist", column.Name)
	}

	t.Columns = append(t.Columns, column)
	t.ColumnNames[column.Name] = len(t.Columns) - 1
	return nil
}

func (t *TableDesc) GetColumn(name string) (*Column, error) {
	index, exists := t.ColumnNames[name]
	if !exists {
		return nil, fmt.Errorf("column %s doesn't exist", name)
	}

	return t.Columns[index], nil
}

func (t *TableDesc) RemoveColumn(name string) (*Column, error) {
	index, exists := t.ColumnNames[name]
	if !exists {
		return nil, fmt.Errorf("column %s doesn't exist", name)
	}

	column := t.Columns[index]
	if column == nil {
		return nil, fmt.Errorf("column %s doesn't exist", name)
	}
	delete(t.ColumnNames, name)
	t.Columns = slices.Delete(t.Columns, index, index+1)
	return column, nil
}

func (t *TableDesc) AddIndex(index *Index) error {
	_, exists := t.IndexNames[index.Name]
	if exists {
		return fmt.Errorf("index %s already exist", index.Name)
	}
	t.Indexes = append(t.Indexes, index)
	t.IndexNames[index.Name] = len(t.Indexes) - 1
	return nil
}

func (t *TableDesc) GetIndex(name string) (*Index, error) {
	index, exists := t.IndexNames[name]
	if !exists {
		return nil, fmt.Errorf("index %s doesn't exist", name)
	}
	return t.Indexes[index], nil
}

func (t *TableDesc) RemoveIndex(name string) (*Index, error) {
	idx, exists := t.IndexNames[name]
	if !exists {
		return nil, fmt.Errorf("index %s doesn't exist", name)
	}

	index := t.Indexes[idx]
	if index == nil {
		return nil, fmt.Errorf("index %s doesn't exist", name)
	}

	delete(t.IndexNames, name)
	t.Indexes = slices.Delete(t.Indexes, idx, idx)
	return index, nil
}

type Table struct {
	Id   common.Id
	Name string
	Desc *TableDesc
}

func NewTable(id common.Id, name string, tableDesc *TableDesc) *Table {
	return &Table{
		Id:   id,
		Name: name,
		Desc: tableDesc,
	}
}
