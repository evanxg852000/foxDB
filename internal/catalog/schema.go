package catalog

import (
	"fmt"
	"sync/atomic"

	"github.com/evanxg852000/foxdb/internal/utils"
)

type Schema struct {
	id           ObjectId
	name         string
	tableNames   map[string]ObjectId
	tables       map[ObjectId]*Table
	nextObjectId atomic.Uint32
}

func NewSchema(oid ObjectId, name string) *Schema {
	return &Schema{
		id:         oid,
		name:       name,
		tableNames: make(map[string]ObjectId),
		tables:     make(map[ObjectId]*Table),
	}
}

func (s *Schema) GetId() ObjectId {
	return s.id
}

func (s *Schema) GetName() string {
	return s.name
}

func (s *Schema) AddTable(name string) (*Table, error) {
	if _, exists := s.tableNames[name]; exists {
		return nil, fmt.Errorf("table %s already exists", name)
	}

	oid := ObjectId(s.nextObjectId.Add(1))
	table := NewTable(oid, name)
	s.tableNames[table.name] = table.id
	s.tables[table.id] = table
	return table, nil
}

func (s *Schema) GetTable(name string) *Table {
	oid, ok := s.tableNames[name]
	if !ok {
		return nil
	}
	table, ok := s.tables[oid]
	utils.Assert(ok, "table id should exist in tables map")
	return table
}

func (s *Schema) RemoveTable(name string) (*Table, error) {
	oid, ok := s.tableNames[name]
	if !ok {
		return nil, fmt.Errorf("table %s does not exist", name)
	}
	table, ok := s.tables[oid]
	utils.Assert(ok, "table id should exist in tables map")
	delete(s.tableNames, name)
	delete(s.tables, oid)
	return table, nil
}

func (s *Schema) ListTables() []*Table {
	tables := make([]*Table, 0, len(s.tables))
	for _, table := range s.tables {
		tables = append(tables, table)
	}
	return tables
}
