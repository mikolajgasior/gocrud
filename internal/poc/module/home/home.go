package modhomeui

import (
	"context"
	"embed"
	"log/slog"
	"net/http"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/handler"
	handhomeui "codeberg.org/mikolajgasior/gocrud/internal/poc/handler/homeui"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/layout"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/module"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

//go:embed html
var embedHTML embed.FS

type UI struct {
	Layout          *layout.Layout
	handlerPageHome http.HandlerFunc
}

func (f *UI) Init(ctx context.Context, input module.InitInput) error {
	logAttrModule := logger.AttrModule(f)

	if f.Layout == nil {
		slog.Error(module.LogLayoutIsNil, logAttrModule)
		return module.LayoutIsNilError
	}

	handlerInstance := handhomeui.New(handhomeui.HandlerInput{
		EmbedHTML: embedHTML,
		Layout:    f.Layout,
	})
	f.handlerPageHome = handlerInstance.PageHome
	return nil
}

func (f *UI) AddHandler(serveMux *http.ServeMux) {
	handler.AddHandlers(serveMux, map[string]http.HandlerFunc{
		"/": f.handlerPageHome,
	})
}

func (f *UI) WrapHandler(wrapper func(http.HandlerFunc) http.HandlerFunc) {
	f.handlerPageHome = wrapper(f.handlerPageHome)
}

func (f *UI) Sitemap() *layout.Sitemap {
	return &layout.Sitemap{
		Pages: map[string]*layout.Page{
			"01": {
				Path:  "/",
				Title: "Home",
				Auth:  layout.AuthorizedAndUnauthorized,
			},
		},
	}
}
