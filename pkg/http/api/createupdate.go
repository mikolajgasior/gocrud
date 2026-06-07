package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/http/jsonreq"
	"codeberg.org/mikolajgasior/gocrud/pkg/http/jsonresp"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

func (h *Handler) handleAPICreateUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request, path, id string) {
	var idInt uint64
	var err error
	var obj interface{}

	pathOpts := h.options.Paths[path]

	if id != "" {
		idInt, err = strconv.ParseUint(id, 10, 64)
		if err != nil {
			jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
				Ok:   true,
				Code: jsonresp.CodeBadRequest,
			})
			return
		}

		// Use the override constructor — skip reading the existing record.
		// The URL id is authoritative; any id in the JSON body is overwritten below.
		obj, err = h.svc.Read(ctx, path, idInt, pathOpts.UpdateConstructor)
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
		if pathOpts.CreateConstructor != nil {
			obj = pathOpts.CreateConstructor()
		} else {
			obj = h.svc.New(path)
		}
	}

	err = jsonreq.Unmarshal(r, &obj)
	if err != nil {
		jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
			Ok:   true,
			Code: jsonresp.CodeUnmarshalRequest,
		})
		return
	}

	// When an update constructor is used the URL id must win over any id in the
	// JSON body, so set it after unmarshalling.
	if id != "" && pathOpts.UpdateConstructor != nil {
		gocrud.ObjSetIDValue(obj, idInt)
	}

	now := time.Now().UTC().Unix()
	userID := uint64(0)
	if h.options.UserIDFunc != nil {
		userID = h.options.UserIDFunc(r)
	}

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
