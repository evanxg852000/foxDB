package catalog

import (
	"sync"
	"sync/atomic"
)

type ObjectId uint32

type RootCatalog struct {
	sync.Mutex
	schemaNames  map[string]ObjectId
	schemas      map[ObjectId]*Schema
	nextObjectId atomic.Uint32
}

func NewRootCatalog() *RootCatalog {
	return &RootCatalog{
		schemaNames: make(map[string]ObjectId),
		schemas:     make(map[ObjectId]*Schema),
	}
}

func (rc *RootCatalog) AddSchema(schema *Schema) {
	schema.id = ObjectId(rc.nextObjectId.Add(1))
	rc.schemaNames[schema.name] = schema.id
	rc.schemas[schema.id] = schema
}

func (rc *RootCatalog) GetSchema(name string) (*Schema, bool) {
	oid, ok := rc.schemaNames[name]
	if !ok {
		return nil, false
	}
	schema, ok := rc.schemas[oid]
	return schema, ok
}

func (rc *RootCatalog) RemoveSchema(name string) *Schema {
	oid, ok := rc.schemaNames[name]
	if !ok {
		return nil
	}
	schema := rc.schemas[oid]
	delete(rc.schemaNames, name)
	delete(rc.schemas, oid)
	return schema
}

func (rc *RootCatalog) ListSchemas() []*Schema {
	schemas := make([]*Schema, 0, len(rc.schemas))
	for _, schema := range rc.schemas {
		schemas = append(schemas, schema)
	}
	return schemas
}
