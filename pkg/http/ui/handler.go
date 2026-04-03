package ui

import (
	"log/slog"
	"net/http"
	"strings"

	"codeberg.org/mikolajgasior/gocrud/pkg/http/urlpath"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	logAttrHandler := logger.AttrHandler(h)
	logAttrPath := logger.AttrPath(r.URL.Path)

	userID := "1"
	userName := "nameless"

	foundPath := ""
	actionPath := ""
	for path := range h.paths {
		if strings.HasPrefix(r.URL.Path, "/"+path+"/") {
			foundPath = path
			actionPath = r.URL.Path[len("/"+foundPath+"/"):]
			slog.Debug("found path", logAttrHandler, logAttrPath, slog.String("path", foundPath), slog.String("actionPath", actionPath))
			continue
		}
	}

	// every path can have a trailing slash but does not have to
	// /create + /form/create (POST only)
	// /read ???
	// /update/1 + /form/update/1 (POST only)
	// /fetch/delete/1 <- GET only (for now)
	// /list + same filters

	if pathCreateRegexp.MatchString(actionPath) && r.Method == http.MethodGet {
		h.handleCreate(r.Context(), foundPath, userID, userName, w, r)
		return
	}

	if pathFormCreateRegexp.MatchString(actionPath) && r.Method == http.MethodPost {
		h.handleFormCreate(r.Context(), foundPath, userID, userName, w, r)
		return
	}

	/*if pathReadRegexp.MatchString(actionPath) {
		id, ok := urlpath.ID(actionPath)
		if !ok {
			redirect(w, h.pathPrefix+"/"+foundPath+"/"+pathPartList)
			return
		}

		h.handleRead(foundPath, id, userID, userName, w, r)
		return
	}*/

	if pathUpdateRegexp.MatchString(actionPath) && r.Method == http.MethodGet {
		id, ok := urlpath.ID(actionPath)
		if !ok {
			redirect(w, h.pathPrefix+"/"+foundPath+"/"+pathPartList)
			return
		}

		h.handleUpdate(r.Context(), foundPath, id, userID, userName, w, r)
		return
	}

	if pathFormUpdateRegexp.MatchString(actionPath) && r.Method == http.MethodPost {
		id, ok := urlpath.ID(actionPath)
		if !ok {
			redirect(w, h.pathPrefix+"/"+foundPath+"/"+pathPartList)
			return
		}

		h.handleFormUpdate(r.Context(), foundPath, id, userID, userName, w, r)
		return
	}

	if pathFetchDeleteRegexp.MatchString(actionPath) && r.Method == http.MethodDelete {
		id, ok := urlpath.ID(actionPath)
		if !ok {
			redirect(w, h.pathPrefix+"/"+foundPath+"/"+pathPartList)
			return
		}

		h.handleFetchDelete(r.Context(), foundPath, id, w, r)
		return
	}

	if pathListRegexp.MatchString(actionPath) {
		h.handleList(r.Context(), foundPath, userID, userName, w, r)
		return
	}

	redirect(w, pathHome)
}
