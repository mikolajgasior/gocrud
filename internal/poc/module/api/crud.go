package modapi

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/handler"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/layout"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/module"
	handcrudapi "codeberg.org/mikolajgasior/gocrud/pkg/http/api"
	"codeberg.org/mikolajgasior/gocrud/pkg/http/cors"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

type API struct {
	Paths   map[string]func() interface{}
	handler http.HandlerFunc
	svc     *svccrud.CRUD
}

func (a *API) Init(ctx context.Context, input module.InitInput) error {
	logAttrModule := logger.AttrModule(a)

	a.svc = svccrud.New(a.Paths, input.DBConn)
	handlerInstance := handcrudapi.New(a.svc, &cors.CORS{}, a.Paths)
	a.handler = handlerInstance.Serve

	pathKeys := make([]string, 0, len(a.Paths))
	for path := range a.Paths {
		pathKeys = append(pathKeys, path)
	}

	if input.CreateTables {
		slog.Info("running crud table creation", slog.String("pages", strings.Join(pathKeys, ",")), logAttrModule)
		err := a.svc.CreateTables(ctx)
		if err != nil {
			slog.Error(module.LogCannotCreateDBTable, logAttrModule, logger.AttrError(err))
			return &module.CreateTableError{
				Err: err,
			}
		}
	}

	return nil
}

func (a *API) AddHandler(serveMux *http.ServeMux) {
	handler.AddHandler(serveMux, "/v0/crudapi", a.handler)
}

func (a *API) Sitemap() *layout.Sitemap {
	return nil
}
