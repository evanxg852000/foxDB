package catalog

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/evanxg852000/foxdb/internal/utils"
)

type ObjectId uint32

type RootCatalog struct {
	sync.RWMutex
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

func (rc *RootCatalog) AddSchema(name string) (*Schema, error) {
	if _, exists := rc.schemaNames[name]; exists {
		return nil, fmt.Errorf("schema %s already exists", name)
	}

	oid := ObjectId(rc.nextObjectId.Add(1))
	schema := NewSchema(oid, name)
	rc.schemaNames[schema.name] = schema.id
	rc.schemas[schema.id] = schema
	return schema, nil
}

func (rc *RootCatalog) GetSchema(name string) *Schema {
	oid, ok := rc.schemaNames[name]
	if !ok {
		return nil
	}
	schema, ok := rc.schemas[oid]
	utils.Assert(ok, "schema id should exist in schemas map")
	return schema
}

func (rc *RootCatalog) RemoveSchema(name string) (*Schema, error) {
	oid, ok := rc.schemaNames[name]
	if !ok {
		return nil, fmt.Errorf("schema %s does not exist", name)
	}
	schema, ok := rc.schemas[oid]
	utils.Assert(ok, "schema id should exist in schemas map")
	delete(rc.schemaNames, name)
	delete(rc.schemas, oid)
	return schema, nil
}

func (rc *RootCatalog) ListSchemas() []*Schema {
	schemas := make([]*Schema, 0, len(rc.schemas))
	for _, schema := range rc.schemas {
		schemas = append(schemas, schema)
	}
	return schemas
}
