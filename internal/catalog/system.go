package catalog

func AddInformationSchema(rootCatalog *RootCatalog) {
	// Add standard information schema, tables, indexes, etc. here as needed.
	infoSchema := NewSchema("information_schema")

	rootCatalog.AddSchema(infoSchema)

}
