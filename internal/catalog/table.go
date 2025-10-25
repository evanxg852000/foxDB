package catalog

import "sync/atomic"

type Table struct {
	id           ObjectId
	name         string
	columnNames  map[string]ObjectId
	columns      map[ObjectId]*Column
	indexNames   map[string]ObjectId
	indexes      map[ObjectId]*Index
	primaryKey   []ObjectId
	nextObjectId atomic.Uint32
}

func NewTable(name string) *Table {
	return &Table{
		name:        name,
		columnNames: make(map[string]ObjectId),
		columns:     make(map[ObjectId]*Column),
		indexNames:  make(map[string]ObjectId),
		indexes:     make(map[ObjectId]*Index),
		primaryKey:  make([]ObjectId, 0),
	}
}

func (t *Table) GetId() ObjectId {
	return t.id
}

func (t *Table) GetName() string {
	return t.name
}

func (t *Table) AddColumn(column *Column) {
	column.id = ObjectId(t.nextObjectId.Add(1))
	t.columnNames[column.name] = column.id
	t.columns[column.id] = column
}

func (t *Table) GetColumn(name string) (*Column, bool) {
	oid, ok := t.columnNames[name]
	if !ok {
		return nil, false
	}
	column, ok := t.columns[oid]
	return column, ok
}

func (t *Table) RemoveColumn(name string) *Column {
	oid, ok := t.columnNames[name]
	if !ok {
		return nil
	}
	column := t.columns[oid]
	delete(t.columnNames, name)
	delete(t.columns, oid)
	return column
}

func (t *Table) ListColumns() []*Column {
	columns := make([]*Column, 0, len(t.columns))
	for _, column := range t.columns {
		columns = append(columns, column)
	}
	return columns
}

func (t *Table) AddIndex(index *Index) {
	index.id = ObjectId(t.nextObjectId.Add(1))
	t.indexNames[index.name] = index.id
	t.indexes[index.id] = index
}

func (t *Table) GetIndex(name string) (*Index, bool) {
	oid, ok := t.indexNames[name]
	if !ok {
		return nil, false
	}
	index, ok := t.indexes[oid]
	return index, ok
}

func (t *Table) RemoveIndex(name string) *Index {
	oid, ok := t.indexNames[name]
	if !ok {
		return nil
	}
	index := t.indexes[oid]
	delete(t.indexNames, name)
	delete(t.indexes, oid)
	return index
}

func (t *Table) ListIndexes() []*Index {
	indexes := make([]*Index, 0, len(t.indexes))
	for _, index := range t.indexes {
		indexes = append(indexes, index)
	}
	return indexes
}
