package physical

import (
	"context"

	"github.com/evanxg852000/foxdb/internal/catalog"
	"github.com/evanxg852000/foxdb/internal/storage"
	"github.com/evanxg852000/foxdb/internal/types"
)

type Scan struct {
	schema *types.DataSchema
}

func NewScan(table *catalog.Table) *Scan {
	columns := make([]types.DataColumn, len(table.ListColumns()))
	for i, col := range table.ListColumns() {
		columns[i] = types.DataColumn{
			Name:     col.GetName(),
			DataType: col.GetDataType(),
		}
	}
	schema := &types.DataSchema{
		Columns: columns,
	}
	return &Scan{
		schema: schema,
	}
}

func (s *Scan) Execute(ctx context.Context, catalog *catalog.RootCatalog, storage *storage.KvStorage) (*types.DataChunk, error) {
	return types.NewChunk(s.schema), nil
}

func (s *Scan) GetSchema() *types.DataSchema {
	return s.schema
}
