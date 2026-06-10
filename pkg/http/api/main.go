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
	CodeForbidden    = "FORBIDDEN"
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

// Route maps a URL path segment to a service registry key and configures
// per-operation behavior. When RegistryKey is empty the URL path segment is
// used as the registry key.
type Route struct {
	RegistryKey    string
	Flags          int64
	AllowedFilters []string // when non-empty, only the listed fields may be used as filters

	// Allow hooks — called with the object and request; return non-nil to reject
	// with 403. AllowUpdate receives the loaded record before the request body is
	// applied, so it reflects the stored state (e.g. the original owner).
	AllowCreate func(obj interface{}, r *http.Request) error
	AllowUpdate func(obj interface{}, r *http.Request) error
	AllowRead   func(obj interface{}, r *http.Request) error
	AllowDelete func(obj interface{}, r *http.Request) error

	// Per-operation constructor overrides. When nil the service's registered
	// constructor for the key is used. Set to use a different struct type for
	// a specific operation — for example a create-only struct with fewer fields
	// that implements a custom insertQueryBuilder.
	CreateConstructor func() interface{}
	UpdateConstructor func() interface{}
	ReadConstructor   func() interface{}
	ListConstructor   func() interface{}
}

type Options struct {
	CORS   cors.CORS
	Routes map[string]Route
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
