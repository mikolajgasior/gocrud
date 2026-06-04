package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/mikolajgasior/gocrud/pkg/http/jsonresp"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

func (h *Handler) handleAPIList(ctx context.Context, w http.ResponseWriter, r *http.Request, path string) {
	params := r.URL.Query()

	limit, _ := strconv.Atoi(params.Get("limit"))
	offset, _ := strconv.Atoi(params.Get("offset"))
	if limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	order := params.Get("order")
	orderDirection := params.Get("order_direction")

	pathOpts := h.options.Paths[path]

	var filterVals, filterOps map[string]string

	if !pathOpts.DisableFilters {
		filterVals = make(map[string]string)
		filterOps = make(map[string]string)

		// Build allowed-filter lookup once (nil map means all fields are allowed).
		var allowed map[string]bool
		if len(pathOpts.AllowedFilters) > 0 {
			allowed = make(map[string]bool, len(pathOpts.AllowedFilters))
			for _, f := range pathOpts.AllowedFilters {
				allowed[f] = true
			}
		}

		for name := range params {
			if filterValRegexp.MatchString(name) {
				field := strings.TrimPrefix(name, filterValPrefix)
				if allowed == nil || allowed[field] {
					filterVals[field] = params.Get(name)
				}
			}
			if filterOpRegexp.MatchString(name) {
				field := strings.TrimPrefix(name, filterOpPrefix)
				if allowed == nil || allowed[field] {
					filterOps[field] = params.Get(name)
				}
			}
		}
	}

	objs, err := h.svc.List(ctx, path, limit, offset, order, orderDirection, filterVals, filterOps, nil, h.options.Paths[path].ListConstructor)
	if err != nil {
		var validErr *svccrud.FilterValidationError
		if errors.As(err, &validErr) {
			message := err.Error()
			jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
				Ok:      true,
				Code:    jsonresp.CodeValidationFailed,
				Message: &message,
			})
			return
		}

		jsonresp.Write(w, http.StatusInternalServerError, &jsonresp.Response{
			Ok:   true,
			Code: CodeServiceError,
		})
		return
	}

	jsonresp.Write(w, http.StatusOK, &jsonresp.Response{
		Ok:   true,
		Code: jsonresp.CodeSuccess,
		Data: objs,
	})
}
