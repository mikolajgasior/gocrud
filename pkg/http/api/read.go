package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"codeberg.org/mikolajgasior/gocrud/pkg/http/jsonresp"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

func (h *Handler) handleAPIRead(ctx context.Context, w http.ResponseWriter, r *http.Request, key string, route Route, id string) {
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
			Ok:   true,
			Code: jsonresp.CodeBadRequest,
		})
		return
	}

	obj, err := h.svc.Read(ctx, key, idInt, route.ReadConstructor)
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
		Data: responseData(obj, h.svc.PasswordFieldNames(key)),
	})
}
