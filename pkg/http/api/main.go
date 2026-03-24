package api

import (
	"regexp"

	"miko.gs/gocrud/pkg/http/cors"
	svccrud "miko.gs/gocrud/pkg/service"
)

const (
	CodeServiceError = "SERVICE_ERROR"
)

const (
	filterValPrefix = "filter_val_"
	filterOpPrefix  = "filter_op_"
)

var (
	filterValRegexp = regexp.MustCompile("^filter_val_[a-zA-Z0-9_]+$")
	filterOpRegexp  = regexp.MustCompile("^filter_op_[a-zA-Z0-9_]+$")
)

type Handler struct {
	cors  *cors.CORS
	paths map[string]func() interface{}
	svc   *svccrud.CRUD
}

func New(svc *svccrud.CRUD, cors *cors.CORS, paths map[string]func() interface{}) *Handler {
	handler := &Handler{
		cors:  cors,
		paths: paths,
		svc:   svc,
	}
	return handler
}
