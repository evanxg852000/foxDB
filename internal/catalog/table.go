package catalog

import (
	"fmt"
	"sync/atomic"

	"github.com/evanxg852000/foxdb/internal/types"
	"github.com/evanxg852000/foxdb/internal/utils"
)

type Table struct {
	id           ObjectId
	name         string
	columnNames  map[string]ObjectId
	columns      map[ObjectId]*Column
	indexNames   map[string]ObjectId
	indexes      map[ObjectId]*Index
	primaryKeys  []ObjectId
	nextObjectId atomic.Uint32
}

func NewTable(oid ObjectId, name string) *Table {
	return &Table{
		id:          oid,
		name:        name,
		columnNames: make(map[string]ObjectId),
		columns:     make(map[ObjectId]*Column),
		indexNames:  make(map[string]ObjectId),
		indexes:     make(map[ObjectId]*Index),
		primaryKeys: make([]ObjectId, 0),
	}
}

func (t *Table) GetId() ObjectId {
	return t.id
}

func (t *Table) GetName() string {
	return t.name
}

func (t *Table) AddColumn(name string, dataType types.DataType, constraints Constraint) (*Column, error) {
	if _, exists := t.columnNames[name]; exists {
		return nil, fmt.Errorf("column %s already exists", name)
	}

	oid := ObjectId(t.nextObjectId.Add(1))
	column := NewColumn(oid, name, dataType, constraints)
	t.columnNames[column.name] = column.id
	t.columns[column.id] = column
	return column, nil
}

func (t *Table) GetColumn(name string) *Column {
	oid, ok := t.columnNames[name]
	if !ok {
		return nil
	}
	column, ok := t.columns[oid]
	utils.Assert(ok, "column id should exist in columns map")
	return column
}

func (t *Table) RemoveColumn(name string) (*Column, error) {
	oid, ok := t.columnNames[name]
	if !ok {
		return nil, fmt.Errorf("column %s does not exist", name)
	}
	column, ok := t.columns[oid]
	utils.Assert(ok, "column id should exist in columns map")
	delete(t.columnNames, name)
	delete(t.columns, oid)
	return column, nil
}

func (t *Table) ListColumns() []*Column {
	columns := make([]*Column, 0, len(t.columns))
	for _, column := range t.columns {
		columns = append(columns, column)
	}
	return columns
}

func (t *Table) AddIndex(name string, columnNames []string, unique bool) (*Index, error) {
	if _, exists := t.indexNames[name]; exists {
		return nil, fmt.Errorf("index %s already exists", name)
	}

	oid := ObjectId(t.nextObjectId.Add(1))
	columnIds := t.columnIdsFromNames(columnNames)
	index := NewIndex(oid, name, columnIds, unique)
	t.indexNames[index.name] = index.id
	t.indexes[index.id] = index
	return index, nil
}

func (t *Table) GetIndex(name string) *Index {
	oid, ok := t.indexNames[name]
	if !ok {
		return nil
	}
	index, ok := t.indexes[oid]
	utils.Assert(ok, "index id should exist in indexes map")
	return index
}

func (t *Table) RemoveIndex(name string) (*Index, error) {
	oid, ok := t.indexNames[name]
	if !ok {
		return nil, fmt.Errorf("index %s does not exist", name)
	}
	index, ok := t.indexes[oid]
	utils.Assert(ok, "index id should exist in indexes map")
	delete(t.indexNames, name)
	delete(t.indexes, oid)
	return index, nil
}

func (t *Table) ListIndexes() []*Index {
	indexes := make([]*Index, 0, len(t.indexes))
	for _, index := range t.indexes {
		indexes = append(indexes, index)
	}
	return indexes
}

func (t *Table) SetPrimaryKeys(columnNames []string) {
	t.primaryKeys = t.columnIdsFromNames(columnNames)
}

func (t *Table) GetPrimaryKeys() []ObjectId {
	return t.primaryKeys
}

func (t *Table) columnIdsFromNames(columnsNames []string) []ObjectId {
	ids := make([]ObjectId, 0, len(columnsNames))
	for _, colName := range columnsNames {
		col := t.GetColumn(colName)
		if col != nil {
			ids = append(ids, col.id)
		}
		//TODO: handle col == nil
	}
	return ids
}
