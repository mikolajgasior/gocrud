package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"miko.gs/struct-crud/pkg/http/jsonresp"
	svccrud "miko.gs/struct-crud/pkg/service"
)

func (h *Handler) handleAPIRead(ctx context.Context, w http.ResponseWriter, r *http.Request, path, id string) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
			Ok:   true,
			Code: jsonresp.CodeBadRequest,
		})
		return
	}

	obj, err := h.svc.Read(ctx, path, idInt)
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
		Data: obj,
	})
}
