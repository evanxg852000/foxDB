package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Record struct {
	tableDesc *DataSchema
	values    []*Value
}

func NewRecord(tableDesc *DataSchema) *Record {
	return &Record{
		tableDesc: tableDesc,
		values:    make([]*Value, len(tableDesc.Columns)),
	}
}

func (r *Record) SetInt(colIndex uint, v int64) error {
	return r.setAt(colIndex, *NewIntValue(v))
}

func (r *Record) SetFloat(colIndex uint, v float64) error {
	return r.setAt(colIndex, *NewFloatValue(v))
}

func (r *Record) SetBool(colIndex uint, v bool) error {
	return r.setAt(colIndex, *NewBoolValue(v))
}

func (r *Record) SetText(colIndex uint, v string) error {
	return r.setAt(colIndex, *NewTextValue(v))
}

func (r *Record) GetInt(colIndex uint) (int64, error) {
	val, err := r.getAt(colIndex)
	if err != nil {
		return 0, err
	}
	return val.Int()
}

func (r *Record) GetFloat(colIndex uint) (float64, error) {
	val, err := r.getAt(colIndex)
	if err != nil {
		return 0, err
	}
	return val.Float()
}

func (r *Record) GetBool(colIndex uint) (bool, error) {
	val, err := r.getAt(colIndex)
	if err != nil {
		return false, err
	}
	return val.Bool()
}

func (r *Record) GetText(colIndex uint) (string, error) {
	val, err := r.getAt(colIndex)
	if err != nil {
		return "", err
	}
	return val.Text()
}

func (r *Record) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, val := range r.values {
		if val == nil {
			binary.Write(buf, binary.LittleEndian, nil)
			continue
		}

		if val.dataType == TYPE_TEXT {
			strVal := val.data.(string)
			err := binary.Write(buf, binary.LittleEndian, uint16(len(strVal)))
			if err != nil {
				return nil, err
			}
			_, err = buf.WriteString(strVal)
			if err != nil {
				return nil, err
			}
			continue
		}

		err := binary.Write(buf, binary.LittleEndian, val.data)
		if err != nil {
			return nil, err
		}

	}
	return buf.Bytes(), nil
}

func (r *Record) Decode(data []byte) error {
	reader := bytes.NewReader(data)
	for idx, column := range r.tableDesc.Columns {
		switch column.DataType {
		case TYPE_INT:
			var v int64
			err := binary.Read(reader, binary.LittleEndian, &v)
			if err != nil {
				return err
			}
			r.SetInt(uint(idx), v)
		case TYPE_FLOAT:
			var v float64
			err := binary.Read(reader, binary.LittleEndian, &v)
			if err != nil {
				return err
			}
			r.SetFloat(uint(idx), v)
		case TYPE_BOOL:
			var v bool
			err := binary.Read(reader, binary.LittleEndian, &v)
			if err != nil {
				return err
			}
			r.SetBool(uint(idx), v)
		case TYPE_TEXT:
			var strLen uint16
			err := binary.Read(reader, binary.LittleEndian, &strLen)
			if err != nil {
				return err
			}

			stringBytes := make([]byte, strLen)
			if strLen > 0 {
				_, err = reader.Read(stringBytes)
				if err != nil {
					return err
				}
			}
			r.SetText(uint(idx), string(stringBytes))
		}
	}
	return nil
}

func (r *Record) setAt(colIndex uint, v Value) error {
	if colIndex >= uint(len(r.values)) {
		return fmt.Errorf("invalid column index: %d", colIndex)
	}

	col := r.tableDesc.Columns[colIndex]
	if col.DataType != v.dataType {
		return fmt.Errorf("")
	}

	r.values[colIndex] = &v
	return nil
}

func (r *Record) getAt(colIndex uint) (*Value, error) {
	if colIndex >= uint(len(r.values)) {
		return nil, fmt.Errorf("invalid column index: %d", colIndex)
	}
	return r.values[colIndex], nil
}
