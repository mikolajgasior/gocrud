package crud

type Options struct {
	TableNamePrefix string
	TagName         string
	Flags           int64
}

const (
	GetCountOnUniq = 1
)
