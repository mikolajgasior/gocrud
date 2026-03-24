package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"miko.gs/gocrud/pkg/http/jsonresp"
	svccrud "miko.gs/gocrud/pkg/service"
)

func (h *Handler) handleAPIDelete(ctx context.Context, w http.ResponseWriter, path, id string) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
			Ok:   true,
			Code: jsonresp.CodeBadRequest,
		})
		return
	}

	err = h.svc.Delete(ctx, path, idInt)
	if err != nil {
		if errors.Is(err, svccrud.NotFoundError) {
			jsonresp.Write(w, http.StatusNotFound, &jsonresp.Response{
				Ok:   true,
				Code: jsonresp.CodeNotFound,
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
	})
}
