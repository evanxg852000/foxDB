package catalog

import "github.com/evanxg852000/foxdb/internal/types"

func AddInformationSchema(rootCatalog *RootCatalog) {
	rootCatalog.Lock()
	defer rootCatalog.Unlock()

	// Add standard information schema, tables, indexes, etc. here as needed.
	infoSchema, _ := rootCatalog.AddSchema("information_schema")

	// schemas/databases
	schemasTable, _ := infoSchema.AddTable("schemas")
	schemasTable.AddColumn("id", types.TYPE_INT, UniqueConstraint)
	schemasTable.AddColumn("name", types.TYPE_TEXT, NoConstraint)
	schemasTable.SetPrimaryKeys([]string{"id"})

	tablesTable, _ := infoSchema.AddTable("tables")
	tablesTable.AddColumn("id", types.TYPE_INT, UniqueConstraint)
	tablesTable.AddColumn("name", types.TYPE_TEXT, NoConstraint)
	tablesTable.AddColumn("schema_id", types.TYPE_INT, NoConstraint)
	tablesTable.AddColumn("sequence_value", types.TYPE_INT, NoConstraint)
	tablesTable.SetPrimaryKeys([]string{"id"})
}
