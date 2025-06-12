package datasources

type DataSourcer interface {
	ReadSource(chan float32) error
}
