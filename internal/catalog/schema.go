package catalog

import "sync/atomic"

type Schema struct {
	id           ObjectId
	name         string
	tableNames   map[string]ObjectId
	tables       map[ObjectId]*Table
	nextObjectId atomic.Uint32
}

func NewSchema(name string) *Schema {
	return &Schema{
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

func (s *Schema) AddTable(table *Table) {
	table.id = ObjectId(s.nextObjectId.Add(1))
	s.tableNames[table.name] = table.id
	s.tables[table.id] = table
}

func (s *Schema) GetTable(name string) (*Table, bool) {
	oid, ok := s.tableNames[name]
	if !ok {
		return nil, false
	}
	table, ok := s.tables[oid]
	return table, ok
}

func (s *Schema) RemoveTable(name string) *Table {
	oid, ok := s.tableNames[name]
	if !ok {
		return nil
	}
	table := s.tables[oid]
	delete(s.tableNames, name)
	delete(s.tables, oid)
	return table
}

func (s *Schema) ListTables() []*Table {
	tables := make([]*Table, 0, len(s.tables))
	for _, table := range s.tables {
		tables = append(tables, table)
	}
	return tables
}
