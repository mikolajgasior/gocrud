package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"miko.gs/struct-crud/pkg/http/jsonreq"
	"miko.gs/struct-crud/pkg/http/jsonresp"
	svccrud "miko.gs/struct-crud/pkg/service"
)

func (h *Handler) handleAPICreateUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request, path, id string) {
	var idInt int64
	var err error

	var obj interface{}

	if id != "" {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
				Ok:   true,
				Code: jsonresp.CodeBadRequest,
			})
			return
		}

		obj, err = h.svc.Read(ctx, path, idInt)
		if err != nil {
			if errors.Is(err, svccrud.NotFoundError) {
				jsonresp.Write(w, http.StatusNotFound, &jsonresp.Response{
					Ok:   true,
					Code: CodeServiceError,
				})
				return
			}
			return
		}
	} else {
		obj = h.svc.New(path)
	}

	err = jsonreq.Unmarshal(r, &obj)
	if err != nil {
		jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
			Ok:   true,
			Code: jsonresp.CodeUnmarshalRequest,
		})
		return
	}

	now := time.Now().UTC().Unix()
	userID := int64(0)

	err = h.svc.Save(ctx, obj, now, userID)
	if err != nil {
		var validErr *svccrud.ModelValidationError
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
	idInt = h.svc.ID(obj)

	status := http.StatusCreated
	code := jsonresp.CodeCreated
	if id != "" {
		status = http.StatusOK
		code = jsonresp.CodeSuccess
	}

	jsonresp.Write(w, status, &jsonresp.Response{
		Ok:   true,
		Code: code,
		Data: map[string]interface{}{
			"id": idInt,
		},
	})
}
