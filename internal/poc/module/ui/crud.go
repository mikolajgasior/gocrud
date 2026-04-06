package modui

import (
	"context"
	"log/slog"
	"net/http"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/handler"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/layout"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/module"
	"codeberg.org/mikolajgasior/gocrud/pkg/http/cors"
	handcrudui "codeberg.org/mikolajgasior/gocrud/pkg/http/ui"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

const (
	uriCrudUi = "/v0/crudui"
)

type UI struct {
	Layout   *layout.Layout
	Paths    map[string]func() interface{}
	XSitemap map[string]map[string]string
	handler  http.HandlerFunc
}

func (u *UI) Init(ctx context.Context, input module.InitInput) error {
	logAttrModule := logger.AttrModule(u)

	if u.Layout == nil {
		slog.Error(module.LogLayoutIsNil, logAttrModule)
		return module.LayoutIsNilError
	}

	svc := svccrud.New(u.Paths, input.DBConn)
	crudUIHandler := handcrudui.New(handcrudui.HandlerInput{
		Svc:        svc,
		CORS:       &cors.CORS{},
		Paths:      u.Paths,
		Layout:     u.Layout,
		PathPrefix: uriCrudUi,
	})
	u.handler = crudUIHandler.Handler

	return nil
}

func (u *UI) AddHandler(serveMux *http.ServeMux) {
	handler.AddHandler(serveMux, uriCrudUi, u.handler)
}

func (u *UI) Sitemap() *layout.Sitemap {
	sitemap := &layout.Sitemap{
		Pages:       make(map[string]*layout.Page, len(u.Paths)),
		XPageGroups: make([]*layout.XPageGroup, 0, len(u.XSitemap)),
	}

	for path := range u.Paths {
		sitemap.Pages[path] = &layout.Page{
			Path:  uriCrudUi + "/" + path + "/list",
			Title: path,
			Auth:  layout.AuthorizedOnly,
		}
	}

	for groupTitle, pages := range u.XSitemap {
		pageGroup := &layout.XPageGroup{
			Title: groupTitle,
			Pages: make([]*layout.Page, 0, len(pages)),
		}

		for pageTitle, pagePath := range pages {
			page := &layout.Page{
				Path:  uriCrudUi + "/" + pagePath + "/list",
				Title: pageTitle,
				Auth:  layout.AuthorizedOnly,
			}
			pageGroup.Pages = append(pageGroup.Pages, page)
		}

		sitemap.XPageGroups = append(sitemap.XPageGroups, pageGroup)
	}

	return sitemap
}
