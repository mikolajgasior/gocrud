package api

import (
	"regexp"

	"codeberg.org/mikolajgasior/gocrud/pkg/http/cors"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
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

type Options struct {
	CORS cors.CORS
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
