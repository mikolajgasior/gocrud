package handhomeui

import (
	"embed"
	"log/slog"
	"net/http"

	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

func (h *Handler) PageHome(w http.ResponseWriter, r *http.Request) {
	logAttrHandler := logger.AttrHandler(h)
	logAttrPath := logger.AttrPath(r.URL.Path)
	uri := "/"

	homeTemplate, _ := embed.FS.ReadFile(h.embedHTML, "html/content_home.html")

	userID := "1"
	userName := "nameless"

	pageBytes, err := h.layout.Render(uri, userID, userName, string(homeTemplate))
	if err != nil {
		slog.Error("error executing template for home page", logAttrHandler, logAttrPath, logger.AttrError(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(pageBytes)
}
