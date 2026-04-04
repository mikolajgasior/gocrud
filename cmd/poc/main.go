package main

import (
	"context"
	"embed"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/app"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/layout"
	modelrestaurant "codeberg.org/mikolajgasior/gocrud/internal/poc/model/restaurant"
	modeltask "codeberg.org/mikolajgasior/gocrud/internal/poc/model/task"
	modelwarehouse "codeberg.org/mikolajgasior/gocrud/internal/poc/model/warehouse"
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

	paths := map[string]func() interface{}{
		// Warehouse Module
		"warehouse/suppliers": func() interface{} {
			return &modelwarehouse.Supplier{}
		},
		"warehouse/products": func() interface{} {
			return &modelwarehouse.Product{}
		},
		"warehouse/categories": func() interface{} {
			return &modelwarehouse.Category{}
		},
		"warehouse/warehouses": func() interface{} {
			return &modelwarehouse.Warehouse{}
		},
		"warehouse/stock_movements": func() interface{} {
			return &modelwarehouse.StockMovement{}
		},
		"warehouse/purchase_orders": func() interface{} {
			return &modelwarehouse.PurchaseOrder{}
		},

		// Restaurant Module
		"restaurant/menu_items": func() interface{} {
			return &modelrestaurant.MenuItem{}
		},
		"restaurant/categories": func() interface{} {
			return &modelrestaurant.Category{}
		},
		"restaurant/tables": func() interface{} {
			return &modelrestaurant.Table{}
		},
		"restaurant/orders": func() interface{} {
			return &modelrestaurant.Order{}
		},
		"restaurant/order_items": func() interface{} {
			return &modelrestaurant.OrderItem{}
		},
		"restaurant/staff": func() interface{} {
			return &modelrestaurant.Staff{}
		},

		// Task Management Module
		"task/projects": func() interface{} {
			return &modeltask.Project{}
		},
		"task/tasks": func() interface{} {
			return &modeltask.Task{}
		},
		"task/users": func() interface{} {
			return &modeltask.User{}
		},
		"task/comments": func() interface{} {
			return &modeltask.Comment{}
		},
		"task/attachments": func() interface{} {
			return &modeltask.Attachment{}
		},
	}

	apiModule := &modapi.API{
		Paths: paths,
	}

	uiModule := &modui.UI{
		Layout: layout,
		Paths:  paths,
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
