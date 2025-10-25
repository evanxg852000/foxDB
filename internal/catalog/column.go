package catalog

import (
	"github.com/evanxg852000/foxdb/internal/types"
)

type Constraint struct {
	Unique  bool
	NotNull bool
}

type Column struct {
	id          ObjectId
	name        string
	dataType    types.DataType
	constraints Constraint
}

func NewColumn(id ObjectId, name string, dataType types.DataType, constraints Constraint) *Column {
	return &Column{
		id:          id,
		name:        name,
		dataType:    dataType,
		constraints: constraints,
	}
}

func (c *Column) GetName() string {
	return c.name
}

func (c *Column) GetDataType() types.DataType {
	return c.dataType
}
