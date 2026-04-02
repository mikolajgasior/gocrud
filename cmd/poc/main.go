package main

import (
	"context"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/app"
	modelpage "codeberg.org/mikolajgasior/gocrud/internal/poc/model/page"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/module"
	modapi "codeberg.org/mikolajgasior/gocrud/internal/poc/module/api"
)

func main() {
	var appObj = app.App{
		Modules: map[string]module.Module{
			"crudapi_pages": &modapi.API{
				Paths: map[string]func() interface{}{
					"pages": func() interface{} {
						return &modelpage.Page{}
					},
					"rows": func() interface{} {
						return &modelpage.Row{}
					},
					"cols": func() interface{} {
						return &modelpage.Col{}
					},
					"elements": func() interface{} {
						return &modelpage.Element{}
					},
				},
			},
		},
	}
	appObj.Run(context.Background())
}
