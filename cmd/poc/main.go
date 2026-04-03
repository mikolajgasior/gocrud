package main

import (
	"context"
	"embed"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/app"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/layout"
	modelpage "codeberg.org/mikolajgasior/gocrud/internal/poc/model/page"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/module"
	modapi "codeberg.org/mikolajgasior/gocrud/internal/poc/module/api"
	modhome "codeberg.org/mikolajgasior/gocrud/internal/poc/module/home"
	modui "codeberg.org/mikolajgasior/gocrud/internal/poc/module/ui"
)

//go:embed html
var embedHTML embed.FS

func main() {

	layout := &layout.Layout{
		HTML: embedHTML,
	}

	apiModule := &modapi.API{
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
	}

	uiModule := &modui.UI{
		Layout: layout,
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
	}

	homeModule := &modhome.UI{
		Layout: layout,
	}

	uiSitemap := uiModule.Sitemap()
	if uiSitemap != nil {
		layout.AddSitemap(uiSitemap)
	}

	homeSitemap := homeModule.Sitemap()
	if homeSitemap != nil {
		layout.AddSitemap(homeSitemap)
	}

	var appObj = app.App{
		Modules: map[string]module.Module{
			"10_api":  apiModule,
			"15_ui":   uiModule,
			"20_home": homeModule,
		},
	}
	appObj.Run(context.Background())
}
