package datasources

type DataSourcer interface {
	ReadSource() (float32, error)
}
