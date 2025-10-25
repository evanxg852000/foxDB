package types

import "fmt"

type Value struct {
	dataType DataType
	data     any
}

func NewIntValue(v int64) *Value {
	return &Value{
		dataType: TYPE_INT,
		data:     v,
	}
}

func NewFloatValue(v float64) *Value {
	return &Value{
		dataType: TYPE_FLOAT,
		data:     v,
	}
}

func NewBoolValue(v bool) *Value {
	return &Value{
		dataType: TYPE_BOOL,
		data:     v,
	}
}

func NewTextValue(v string) *Value {
	return &Value{
		dataType: TYPE_TEXT,
		data:     v,
	}
}

func (v *Value) Int() (int64, error) {
	if v.dataType != TYPE_INT {
		return 0, fmt.Errorf("value is not of type INT")
	}
	return v.data.(int64), nil
}

func (v *Value) Float() (float64, error) {
	if v.dataType != TYPE_FLOAT {
		return 0, fmt.Errorf("value is not of type FLOAT")
	}
	return v.data.(float64), nil
}

func (v *Value) Bool() (bool, error) {
	if v.dataType != TYPE_BOOL {
		return false, fmt.Errorf("value is not of type BOOL")
	}
	return v.data.(bool), nil
}

func (v *Value) Text() (string, error) {
	if v.dataType != TYPE_TEXT {
		return "", fmt.Errorf("value is not of type TEXT")
	}
	return v.data.(string), nil
}
