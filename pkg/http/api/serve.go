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
