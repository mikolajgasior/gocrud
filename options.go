package crud

type Options struct {
	TableNamePrefix string
	TagName         string
	Flags           uint64
}

const (
	GetCountOnUniq = 1
)
