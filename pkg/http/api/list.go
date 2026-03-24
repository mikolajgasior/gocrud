package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"miko.gs/gocrud/pkg/http/jsonresp"
	svccrud "miko.gs/gocrud/pkg/service"
)

func (h *Handler) handleAPIList(ctx context.Context, w http.ResponseWriter, r *http.Request, path string) {
	// TODO: not all fields should be filterable

	params := r.URL.Query()

	names := make([]string, 0, len(params))
	for name := range params {
		names = append(names, name)
	}

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

	filterVals := make(map[string]string, len(params))
	filterOps := make(map[string]string, len(params))
	for _, name := range names {
		if filterValRegexp.MatchString(name) {
			filterVals[strings.Replace(name, filterValPrefix, "", 1)] = params.Get(name)
		}
		if filterOpRegexp.MatchString(name) {
			filterOps[strings.Replace(name, filterOpPrefix, "", 1)] = params.Get(name)
		}
	}

	objs, err := h.svc.List(ctx, path, limit, offset, order, orderDirection, filterVals, filterOps, nil)
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
