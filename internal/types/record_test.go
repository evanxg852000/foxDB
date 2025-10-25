package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordEncodeDecodeInt(t *testing.T) {
	// Create table descriptor with int column
	tableDesc := &DataSchema{
		Columns: []DataColumn{
			{Name: "columnName", DataType: TYPE_INT},
		},
	}

	// Create record and set int value
	record := NewRecord(tableDesc)
	err := record.SetInt(0, 12345)
	require.NoError(t, err)

	// Encode
	encoded, err := record.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	// Create new record and decode
	decodedRecord := NewRecord(tableDesc)
	err = decodedRecord.Decode(encoded)
	require.NoError(t, err)

	// Verify decoded value
	value, err := decodedRecord.GetInt(0)
	require.NoError(t, err)
	assert.Equal(t, int64(12345), value)
}

func TestRecordEncodeDecodeFloat(t *testing.T) {
	// Create table descriptor with float column
	tableDesc := &DataSchema{
		Columns: []DataColumn{
			{Name: "columnName", DataType: TYPE_FLOAT},
		},
	}

	// Create record and set float value
	record := NewRecord(tableDesc)
	err := record.SetFloat(0, 3.14159)
	require.NoError(t, err)

	// Encode
	encoded, err := record.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	// Create new record and decode
	decodedRecord := NewRecord(tableDesc)
	err = decodedRecord.Decode(encoded)
	require.NoError(t, err)

	// Verify decoded value
	value, err := decodedRecord.GetFloat(0)
	require.NoError(t, err)
	assert.Equal(t, 3.14159, value)
}

func TestRecordEncodeDecodeBool(t *testing.T) {
	// Create table descriptor with bool column
	tableDesc := &DataSchema{
		Columns: []DataColumn{
			{Name: "columnName", DataType: TYPE_BOOL},
		},
	}

	// Create record and set bool value
	record := NewRecord(tableDesc)
	err := record.SetBool(0, true)
	require.NoError(t, err)

	// Encode
	encoded, err := record.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	// Create new record and decode
	decodedRecord := NewRecord(tableDesc)
	err = decodedRecord.Decode(encoded)
	require.NoError(t, err)

	// Verify decoded value
	value, err := decodedRecord.GetBool(0)
	require.NoError(t, err)
	assert.Equal(t, true, value)
}

func TestRecordEncodeDecodeText(t *testing.T) {
	// Create table descriptor with text column
	tableDesc := &DataSchema{
		Columns: []DataColumn{
			{Name: "columnName", DataType: TYPE_TEXT},
		},
	}

	// Create record and set text value
	record := NewRecord(tableDesc)
	err := record.SetText(0, "hello world")
	require.NoError(t, err)

	// Encode
	encoded, err := record.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	// Create new record and decode
	decodedRecord := NewRecord(tableDesc)
	err = decodedRecord.Decode(encoded)
	require.NoError(t, err)

	// Verify decoded value
	value, err := decodedRecord.GetText(0)
	require.NoError(t, err)
	assert.Equal(t, "hello world", value)
}

func TestRecordEncodeDecodeMultipleColumns(t *testing.T) {
	// Create table descriptor with multiple columns
	tableDesc := &DataSchema{
		Columns: []DataColumn{
			{Name: "columnName1", DataType: TYPE_INT},
			{Name: "columnName2", DataType: TYPE_FLOAT},
			{Name: "columnName3", DataType: TYPE_BOOL},
			{Name: "columnName4", DataType: TYPE_TEXT},
		},
	}

	// Create record and set values
	record := NewRecord(tableDesc)
	err := record.SetInt(0, 42)
	require.NoError(t, err)
	err = record.SetFloat(1, 2.71828)
	require.NoError(t, err)
	err = record.SetBool(2, false)
	require.NoError(t, err)
	err = record.SetText(3, "test string")
	require.NoError(t, err)

	// Encode
	encoded, err := record.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	// Create new record and decode
	decodedRecord := NewRecord(tableDesc)
	err = decodedRecord.Decode(encoded)
	require.NoError(t, err)

	// Verify all decoded values
	intVal, err := decodedRecord.GetInt(0)
	require.NoError(t, err)
	assert.Equal(t, int64(42), intVal)

	floatVal, err := decodedRecord.GetFloat(1)
	require.NoError(t, err)
	assert.Equal(t, 2.71828, floatVal)

	boolVal, err := decodedRecord.GetBool(2)
	require.NoError(t, err)
	assert.Equal(t, false, boolVal)

	textVal, err := decodedRecord.GetText(3)
	require.NoError(t, err)
	assert.Equal(t, "test string", textVal)
}

func TestRecordEncodeDecodeEmptyString(t *testing.T) {
	// Create table descriptor with text column
	tableDesc := &DataSchema{
		Columns: []DataColumn{
			{Name: "columnName", DataType: TYPE_TEXT},
		},
	}

	// Create record and set empty string
	record := NewRecord(tableDesc)
	err := record.SetText(0, "")
	require.NoError(t, err)

	// Encode
	encoded, err := record.Encode()
	require.NoError(t, err)

	// Create new record and decode
	decodedRecord := NewRecord(tableDesc)
	err = decodedRecord.Decode(encoded)
	require.NoError(t, err)

	// Verify decoded value
	value, err := decodedRecord.GetText(0)
	require.NoError(t, err)
	assert.Equal(t, "", value)
}
