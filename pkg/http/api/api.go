package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"miko.gs/struct-crud/pkg/http/cors"
	"miko.gs/struct-crud/pkg/http/jsonreq"
	"miko.gs/struct-crud/pkg/http/jsonresp"
	"miko.gs/struct-crud/pkg/http/urlpath"
	"miko.gs/struct-crud/pkg/logger"
	svccrud "miko.gs/struct-crud/pkg/service"
)

const (
	CodeServiceError = "SERVICE_ERROR"
)

const (
	filterValPrefix = "filter_val_"
	filterOpPrefix  = "filter_op_"
)

var (
	filterValRegexp = regexp.MustCompile("^filter_val_[a-zA-Z0-9_]+$")
	filterOpRegexp  = regexp.MustCompile("^filter_op_[a-zA-Z0-9_]+$")
)

type Handler struct {
	cors  *cors.CORS
	paths map[string]func() interface{}
	svc   *svccrud.CRUD
}

func New(svc *svccrud.CRUD, cors *cors.CORS, paths map[string]func() interface{}) *Handler {
	handler := &Handler{
		cors:  cors,
		paths: paths,
		svc:   svc,
	}
	return handler
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	logAttrHandler := logger.AttrHandler(h)
	logAttrPath := logger.AttrPath(r.URL.Path)

	id := ""
	foundPath := ""
	for path := range h.paths {
		if strings.HasPrefix(r.URL.Path, "/"+path+"/") {
			foundPath = path
			continue
		}
	}

	if foundPath == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	idPath := r.URL.Path[len("/"+foundPath+"/"):]
	id, ok := urlpath.ID(idPath)
	if !ok {
		slog.Debug("ID in URI not numeric", logAttrHandler, logAttrPath)
		jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
			Ok:   true,
			Code: jsonresp.CodeURLPathID,
		})
		return
	}

	// create and update
	if r.Method == http.MethodPut {
		h.handleAPICreateUpdate(r.Context(), w, r, foundPath, id)
		return
	}

	// delete
	if r.Method == http.MethodDelete && id != "" {
		h.handleAPIDelete(r.Context(), w, foundPath, id)
		return
	}

	// read
	if r.Method == http.MethodGet && id != "" {
		h.handleAPIRead(r.Context(), w, r, foundPath, id)
		return
	}

	if r.Method == http.MethodGet && id == "" {
		h.handleAPIList(r.Context(), w, r, foundPath)
		return
	}

	jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
		Ok:   true,
		Code: jsonresp.CodeBadRequest,
	})
	return
}

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

func (h *Handler) handleAPIList(ctx context.Context, w http.ResponseWriter, r *http.Request, path string) {
	// TODO: not all fields should be filterable

	params := r.URL.Query()

	names := make([]string, 0, len(params))
	for name := range params {
		names = append(names, name)
	}

	limit, _ := strconv.Atoi(params.Get("limit"))
	offset, _ := strconv.Atoi(params.Get("offset"))
	if limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	order := params.Get("order")
	orderDirection := params.Get("order_direction")

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

	objs, err := h.svc.List(ctx, path, limit, offset, order, orderDirection, filterVals, filterOps, nil)
	if err != nil {
		var validErr *svccrud.FilterValidationError
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

	jsonresp.Write(w, http.StatusOK, &jsonresp.Response{
		Ok:   true,
		Code: jsonresp.CodeSuccess,
		Data: objs,
	})
}
