package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"codeberg.org/mikolajgasior/gocrud/pkg/http/jsonresp"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

func (h *Handler) handleAPIDelete(ctx context.Context, w http.ResponseWriter, key, id string) {
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
			Ok:   true,
			Code: jsonresp.CodeBadRequest,
		})
		return
	}

	err = h.svc.Delete(ctx, key, idInt)
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
