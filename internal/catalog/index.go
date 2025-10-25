package catalog

type Index struct {
	id        ObjectId
	name      string
	columnIds []ObjectId
	unique    bool
}
