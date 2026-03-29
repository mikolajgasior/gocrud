package ui

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

func (h *Handler) handleFormUpdate(ctx context.Context, path string, id string, userID, userName string, w http.ResponseWriter, r *http.Request) {
	logAttrHandler := logger.AttrHandler(h)
	logAttrPath := logger.AttrPath(r.URL.Path)

	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		errCode := logger.LogError("error converting object id to int", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, h.pathPrefix+"/"+path+"/"+pathPartList+"?error=unknown&error_code="+errCode)
		return
	}

	obj, err := h.svc.Read(ctx, path, idInt)
	if err != nil {
		if errors.Is(err, svccrud.NotFoundError) {
			redirect(w, h.pathPrefix+"/"+path+"/"+pathPartList+"?error=not_found")
			return
		}

		errCode := logger.LogError("error with service reading an object", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, h.pathPrefix+"/"+path+"/"+pathPartList+"?error=unknown&error_code="+errCode)
		return
	}

	now := uint64(time.Now().UTC().Unix())
	userIDInt, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errCode := logger.LogError("error converting user id to int", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, h.pathPrefix+"/"+path+"/"+pathPartUpdate+"/"+id+"?error=unauthorized&error_code="+errCode)
		return
	}

	err = r.ParseForm()
	if err != nil {
		errCode := logger.LogError("error parsing form", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, h.pathPrefix+"/"+path+"/"+pathPartUpdate+"/"+id+"?error=invalid_form_data&error_code="+errCode)
		return
	}

	err = h.svc.SaveFromForm(ctx, obj, r.Form, htmlInputNamePrefix, now, userIDInt)
	if err != nil {
		var validErr *svccrud.ModelValidationError
		if errors.As(err, &validErr) {
			violatedFields := make([]string, len(validErr.Violations))
			urlErrorFields := ""
			for violatedField, _ := range validErr.Violations {
				urlErrorFields += "&error_field=" + violatedField
				violatedFields = append(violatedFields, violatedField)
			}
			violationsString := strings.Join(violatedFields, ", ")
			slog.Error("error validating object for update", logAttrHandler, logAttrPath, logger.AttrError(err), slog.String("violations", violationsString))
			redirect(w, h.pathPrefix+"/"+path+"/"+pathPartUpdate+"/"+id+"?error=invalid_form_data"+urlErrorFields)
			return
		}

		errCode := logger.LogError("error with service saving from form", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, pathHome+"?error=unknown&error_code="+errCode)
		return
	}

	redirect(w, h.pathPrefix+"/"+path+"/"+pathPartUpdate+"/"+id+"?status=updated")
}
