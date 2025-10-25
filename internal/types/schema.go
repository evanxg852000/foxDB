package types

type DataColumn struct {
	Name     string
	DataType DataType
}

type DataSchema struct {
	Columns []DataColumn
}

type DataRow struct {
	Values []Value
}

// import (
// 	"fmt"
// 	"slices"
// )

// type Schema struct {
// 	Tables     []*Table       `json:"-"`
// 	TableNames map[string]int `json:"table_names"`
// }

// func NewSchema() *Schema {
// 	return &Schema{
// 		Tables:     make([]*Table, 0),
// 		TableNames: make(map[string]int),
// 	}
// }

// func (s *Schema) AddTable(table *Table) error {
// 	_, exists := s.TableNames[table.Name]
// 	if exists {
// 		return fmt.Errorf("table %s already exist", table.Name)
// 	}

// 	s.Tables = append(s.Tables, table)
// 	s.TableNames[table.Name] = len(s.Tables) - 1
// 	return nil
// }

// func (s *Schema) GetTable(name string) (*Table, error) {
// 	index, exists := s.TableNames[name]
// 	if !exists {
// 		return nil, fmt.Errorf("table %s doesn't exist", name)
// 	}

// 	return s.Tables[index], nil
// }

// func (s *Schema) RemoveTable(name string) (*Table, error) {
// 	index, exists := s.TableNames[name]
// 	if !exists {
// 		return nil, fmt.Errorf("table %s doesn't exist", name)
// 	}

// 	table := s.Tables[index]
// 	if table == nil {
// 		return nil, fmt.Errorf("table %s doesn't exist", name)
// 	}
// 	delete(s.TableNames, name)
// 	s.Tables = slices.Delete(s.Tables, index, index+1)
// 	return table, nil
// }
