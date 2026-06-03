package crud

type Options struct {
	TableNamePrefix string
	TagName         string
	Dialect         string
	Flags           uint64
}

const (
	GetCountOnUniq = 1
)

const (
	DialectPostgres = "postgres"
	DialectSQLite   = "sqlite"
)
