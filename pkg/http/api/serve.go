package api

import (
	"log/slog"
	"net/http"
	"strings"

	"codeberg.org/mikolajgasior/gocrud/pkg/http/jsonresp"
	"codeberg.org/mikolajgasior/gocrud/pkg/http/urlpath"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

func (h *Handler) Serve(w http.ResponseWriter, r *http.Request) {
	h.options.CORS.WriteHeaders(w)

	logAttrHandler := logger.AttrHandler(h)
	logAttrPath := logger.AttrPath(r.URL.Path)

	foundPath := ""
	for _, path := range h.svc.Paths() {
		if strings.HasPrefix(r.URL.Path, "/"+path+"/") {
			foundPath = path
			break
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

	pathOpts := h.options.Paths[foundPath]

	// create and update
	if r.Method == http.MethodPut {
		if id == "" && pathOpts.DisableCreate {
			jsonresp.Write(w, http.StatusMethodNotAllowed, &jsonresp.Response{
				Ok:   true,
				Code: CodeNotAllowed,
			})
			return
		}
		if id != "" && pathOpts.DisableUpdate {
			jsonresp.Write(w, http.StatusMethodNotAllowed, &jsonresp.Response{
				Ok:   true,
				Code: CodeNotAllowed,
			})
			return
		}
		h.handleAPICreateUpdate(r.Context(), w, r, foundPath, id)
		return
	}

	// delete
	if r.Method == http.MethodDelete && id != "" {
		if pathOpts.DisableDelete {
			jsonresp.Write(w, http.StatusMethodNotAllowed, &jsonresp.Response{
				Ok:   true,
				Code: CodeNotAllowed,
			})
			return
		}
		h.handleAPIDelete(r.Context(), w, foundPath, id)
		return
	}

	// read
	if r.Method == http.MethodGet && id != "" {
		if pathOpts.DisableRead {
			jsonresp.Write(w, http.StatusMethodNotAllowed, &jsonresp.Response{
				Ok:   true,
				Code: CodeNotAllowed,
			})
			return
		}
		h.handleAPIRead(r.Context(), w, r, foundPath, id)
		return
	}

	// list
	if r.Method == http.MethodGet && id == "" {
		if pathOpts.DisableList {
			jsonresp.Write(w, http.StatusMethodNotAllowed, &jsonresp.Response{
				Ok:   true,
				Code: CodeNotAllowed,
			})
			return
		}
		h.handleAPIList(r.Context(), w, r, foundPath)
		return
	}

	jsonresp.Write(w, http.StatusBadRequest, &jsonresp.Response{
		Ok:   true,
		Code: jsonresp.CodeBadRequest,
	})
}
