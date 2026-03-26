package ui

import (
	"bytes"
	"context"
	"embed"
	"errors"
	htmltemplate "html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

func (h *Handler) handleList(ctx context.Context, path string, userID, userName string, w http.ResponseWriter, r *http.Request) {
	logAttrHandler := logger.AttrHandler(h)
	logAttrPath := logger.AttrPath(r.URL.Path)

	params := r.URL.Query()

	names := make([]string, 0, len(params))
	for name := range params {
		names = append(names, name)
	}

	limit, _ := strconv.Atoi(params.Get("limit"))
	page, _ := strconv.Atoi(params.Get("page"))
	if limit < 1 {
		limit = 10
	}
	if page < 0 {
		page = 1
	}

	offset := (page - 1) * limit

	order := params.Get("order")
	orderDirection := params.Get("order_direction")

	if order == "" {
		order = "ID"
	}

	if orderDirection == "" {
		orderDirection = "ASC"
	}

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

	obj := h.svc.New(path)
	fieldNames := structFieldNames(obj)

	objs, err := h.svc.List(ctx, path, limit, offset, order, orderDirection, filterVals, filterOps, func(obj interface{}) interface{} {
		return structToMap(obj)
	})
	if err != nil {
		var validErr *svccrud.FilterValidationError
		if errors.As(err, &validErr) {
			violatedFilters := make([]string, len(validErr.Violations))
			urlErrorFilters := ""
			for violatedFilter, _ := range validErr.Violations {
				urlErrorFilters += "&error_filter=" + violatedFilter
				violatedFilters = append(violatedFilters, violatedFilter)
			}
			violationsString := strings.Join(violatedFilters, ", ")

			slog.Error("error validating filters for listing", logAttrHandler, logAttrPath, logger.AttrError(err), slog.String("violations", violationsString))
			redirect(w, h.pathPrefix+"/"+path+"/"+pathPartList+"?error=invalid_filters"+urlErrorFilters)
			return
		}

		errCode := logger.LogError("error with service getting list of objs", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, pathHome+"?error=unknown&error_code="+errCode)
		return
	}

	numObjs, err := h.svc.Num(ctx, path, filterVals, filterOps)
	if err != nil {
		// Filters are the same as above so we're assuming that they are valid
		errCode := logger.LogError("error with service getting num of objs", logAttrHandler, logAttrPath, logger.AttrError(err))
		redirect(w, pathHome+"?error=unknown&error_code="+errCode)
		return
	}

	pages := int(numObjs) / limit
	if int(numObjs)%limit != 0 {
		pages++
	}

	objectTemplate, _ := embed.FS.ReadFile(h.embedHTML, "html/content_list.html")

	tplObj := struct {
		ObjectCreatePath      string
		ObjectUpdatePath      string
		ObjectFetchDeletePath string
		ObjectListPath        string
		Objects               []interface{}
		NumObjects            int64
		NumPages              int
		ObjectFieldNames      []string
		ParamLimit            int
		ParamPage             int
		ParamOrder            string
		ParamOrderDirection   string
		ParamFilterVals       map[string]string
		ParamFilterOps        map[string]string
	}{
		ObjectCreatePath:      h.pathPrefix + "/" + path + "/create",
		ObjectUpdatePath:      h.pathPrefix + "/" + path + "/" + pathPartUpdate,
		ObjectFetchDeletePath: h.pathPrefix + "/" + path + "/fetch/" + pathPartDelete,
		ObjectListPath:        h.pathPrefix + "/" + path + "/" + pathPartList,
		Objects:               objs,
		NumObjects:            numObjs,
		NumPages:              pages,
		ObjectFieldNames:      fieldNames,
		ParamLimit:            limit,
		ParamPage:             page,
		ParamOrder:            order,
		ParamOrderDirection:   orderDirection,
		ParamFilterVals:       filterVals,
		ParamFilterOps:        filterOps,
	}

	buf := &bytes.Buffer{}
	t := htmltemplate.Must(htmltemplate.New("login").Funcs(htmltemplate.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}).Parse(string(objectTemplate)))
	err = t.Execute(buf, &tplObj)
	if err != nil {
		errCode := logger.LogError("error executing template for object list", logAttrHandler, logAttrPath, logger.AttrError(err))
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
