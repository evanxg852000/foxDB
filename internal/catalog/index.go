package catalog

type Index struct {
	id        ObjectId
	name      string
	columnIds []ObjectId
	unique    bool
}

func NewIndex(id ObjectId, name string, columnIds []ObjectId, unique bool) *Index {
	return &Index{
		id:        id,
		name:      name,
		columnIds: columnIds,
		unique:    unique,
	}
}

func (idx *Index) GetName() string {
	return idx.name
}

func (idx *Index) IsUnique() bool {
	return idx.unique
}

func (idx *Index) GetColumnIds() []ObjectId {
	return idx.columnIds
}
