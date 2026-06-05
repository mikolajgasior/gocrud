package ui

import (
	"bytes"
	"context"
	"embed"
	"errors"
	htmltemplate "html/template"
	"net/http"
	"strconv"

	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
	uiform "codeberg.org/mikolajgasior/gocrud/pkg/ui/form"
	uiinput "codeberg.org/mikolajgasior/gocrud/pkg/ui/input"
)

func (h *Handler) handleUpdate(ctx context.Context, path string, id string, userID, userName string, w http.ResponseWriter, r *http.Request) {
	logAttrHandler := logger.AttrHandler(h)
	logAttrPath := logger.AttrPath(r.URL.Path)

	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		errCode := logger.LogError("error converting object id to int", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, h.pathPrefix+"/"+path+"/"+pathPartList+"?error=unknown&error_code="+errCode)
		return
	}

	obj, err := h.svc.Read(ctx, path, idInt, nil)
	if err != nil {
		if errors.Is(err, svccrud.NotFoundError) {
			redirect(w, h.pathPrefix+"/"+path+"/"+pathPartList+"?error=not_found")
			return
		}

		errCode := logger.LogError("error with service reading an object", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, h.pathPrefix+"/"+path+"/"+pathPartList+"?error=unknown&error_code="+errCode)
		return
	}
	objTypeName := structName(obj)

	objectTemplate, _ := embed.FS.ReadFile(h.embedHTML, "html/content_update.html")

	form, err := uiform.BuildForm(obj, &uiform.FormOptions{
		Options: uiinput.Options{
			IDPrefix:   htmlInputIdPrefix,
			NamePrefix: htmlInputNamePrefix,
			ExcludeFields: map[string]struct{}{
				"ID":         struct{}{},
				"CreatedAt":  struct{}{},
				"ModifiedAt": struct{}{},
				"CreatedBy":  struct{}{},
				"ModifiedBy": struct{}{},
			},
			TagName: "crud",
			Values:  true,
		},
		Path: h.pathPrefix + "/" + path + "/form/" + pathPartUpdate + "/" + id,
	})
	if err != nil {
		errCode := logger.LogError("error building form", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, pathHome+"?error=unknown&error_code="+errCode)
		return
	}

	formHTML := form.HTML()

	urlErrors := r.URL.Query()["error"]
	urlErrorCodes := r.URL.Query()["error_code"]
	urlStatuses := r.URL.Query()["status"]
	urlWarnings := r.URL.Query()["warning"]
	urlErrorFields := r.URL.Query()["error_field"]

	tplObj := struct {
		FormHTML       htmltemplate.HTML
		Errors         []string
		ErrorCodes     []string
		Statuses       []string
		Warnings       []string
		ErrorFields    []string
		ObjectTypeName string
	}{
		FormHTML:       formHTML,
		Errors:         urlErrors,
		ErrorCodes:     urlErrorCodes,
		Statuses:       urlStatuses,
		Warnings:       urlWarnings,
		ErrorFields:    urlErrorFields,
		ObjectTypeName: objTypeName,
	}

	buf := &bytes.Buffer{}
	t := htmltemplate.Must(htmltemplate.New("update").Parse(string(objectTemplate)))
	err = t.Execute(buf, &tplObj)
	if err != nil {
		errCode := logger.LogError("error executing template for object update", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, pathHome+"?error=unknown&error_code="+errCode)
		return
	}

	pageBytes, err := h.layout.Render("/"+path+"/", userID, userName, string(buf.Bytes()))
	if err != nil {
		errCode := logger.LogError("error executing template for object page", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, pathHome+"?error=unknown&error_code="+errCode)
		return
	}

	_, _ = w.Write(pageBytes)
}
