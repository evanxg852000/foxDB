package physical

import (
	"context"

	"github.com/evanxg852000/foxdb/internal/catalog"
	"github.com/evanxg852000/foxdb/internal/types"
)

type ScanExec struct {
	schema *types.DataSchema
}

func NewScanExec(table *catalog.Table) *ScanExec {
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
	return &ScanExec{
		schema: schema,
	}
}

func (s *ScanExec) Execute(ctx context.Context) (*types.DataChunk, error) {
	return types.NewChunk(s.schema), nil
}

func (s *ScanExec) GetSchema() *types.DataSchema {
	return s.schema
}
