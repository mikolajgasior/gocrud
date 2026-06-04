package api

import (
	"regexp"

	"codeberg.org/mikolajgasior/gocrud/pkg/http/cors"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

const (
	CodeServiceError = "SERVICE_ERROR"
	CodeNotAllowed   = "NOT_ALLOWED"
)

const (
	filterValPrefix = "filter_val_"
	filterOpPrefix  = "filter_op_"
)

var (
	filterValRegexp = regexp.MustCompile("^filter_val_[a-zA-Z0-9_]+$")
	filterOpRegexp  = regexp.MustCompile("^filter_op_[a-zA-Z0-9_]+$")
)

// PathOptions controls which operations and filters are available for a single path.
// The zero value enables everything — only set fields you want to restrict.
type PathOptions struct {
	DisableCreate  bool
	DisableUpdate  bool
	DisableDelete  bool
	DisableRead    bool
	DisableList    bool
	DisableFilters bool     // when true, all filter query parameters are ignored
	AllowedFilters []string // when non-empty, only the listed fields may be used as filters
}

type Options struct {
	CORS  cors.CORS
	Paths map[string]PathOptions
}

type Handler struct {
	options Options
	svc     *svccrud.CRUD
}

func New(svc *svccrud.CRUD, options Options) *Handler {
	return &Handler{
		options: options,
		svc:     svc,
	}
}
