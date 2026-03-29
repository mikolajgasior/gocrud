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

func (h *Handler) handleFormCreate(ctx context.Context, path string, userID, userName string, w http.ResponseWriter, r *http.Request) {
	logAttrHandler := logger.AttrHandler(h)
	logAttrPath := logger.AttrPath(r.URL.Path)

	obj := h.svc.New(path)

	now := uint64(time.Now().UTC().Unix())
	userIDInt, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errCode := logger.LogError("error converting user id to int", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, h.pathPrefix+"/"+path+"/"+pathPartCreate+"?error=unauthorized&error_code="+errCode)
		return
	}

	err = r.ParseForm()
	if err != nil {
		errCode := logger.LogError("error parsing form", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, h.pathPrefix+"/"+path+"/"+pathPartCreate+"?error=invalid_form_data&error_code="+errCode)
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
			slog.Error("error validating object for creation", logAttrHandler, logAttrPath, logger.AttrError(err), slog.String("violations", violationsString))
			redirect(w, h.pathPrefix+"/"+path+"/"+pathPartCreate+"?error=invalid_form_data"+urlErrorFields)
			return
		}

		errCode := logger.LogError("error with service saving from form", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, pathHome+"?error=unknown&error_code="+errCode)
		return
	}

	redirect(w, h.pathPrefix+"/"+path+"/"+pathPartCreate+"?status=created")
}
