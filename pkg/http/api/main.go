package api

import (
	"net/http"
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

const (
	DisableCreate = 1 << iota
	DisableUpdate
	DisableDelete
	DisableRead
	DisableList
	DisableFilters
)

var (
	filterValRegexp = regexp.MustCompile("^filter_val_[a-zA-Z0-9_]+$")
	filterOpRegexp  = regexp.MustCompile("^filter_op_[a-zA-Z0-9_]+$")
)

// PathOptions controls which operations and filters are available for a single
// path. The zero value enables everything — only set fields you want to restrict.
type PathOptions struct {
	Flags          int64
	AllowedFilters []string // when non-empty, only the listed fields may be used as filters

	// Per-operation constructor overrides. When nil the service's registered
	// constructor for the path is used. Set to use a different struct type for
	// a specific operation — for example a create-only struct with fewer fields
	// that implements a custom insertQueryBuilder.
	CreateConstructor func() interface{}
	UpdateConstructor func() interface{}
	ReadConstructor   func() interface{}
	ListConstructor   func() interface{}
}

type Options struct {
	CORS  cors.CORS
	Paths map[string]PathOptions
	// UserIDFunc is called on every create and update request to obtain the
	// current user's ID, which is passed to the service as ModifiedBy.
	// When nil, ModifiedBy is always 0.
	UserIDFunc func(r *http.Request) uint64
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
