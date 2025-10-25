package types

type DataChunk struct {
	schema *DataSchema
	rows   []DataRow
}

func NewChunk(schema *DataSchema) *DataChunk {
	return &DataChunk{
		schema: schema,
		rows:   make([]DataRow, 0),
	}
}

func NewWith(schema *DataSchema, rows []DataRow) *DataChunk {
	return &DataChunk{schema, rows}
}

func (c *DataChunk) GetSchema() *DataSchema {
	return c.schema
}

func (c *DataChunk) GetRows() []DataRow {
	return c.rows
}

func (c *DataChunk) AppendRow(row DataRow) {
	c.rows = append(c.rows, row)
}

func (c *DataChunk) GetColumnNames() []string {
	names := make([]string, len(c.schema.Columns))
	for i, col := range c.schema.Columns {
		names[i] = col.Name
	}
	return names
}
