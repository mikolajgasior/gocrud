package main

import (
	"context"

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

func main() {

	layoutInstance := &layout.Layout{}

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
		Layout: layoutInstance,
		Paths:  paths,
		XSitemap: map[string]map[string]string{
			"Warehouse": {
				"Suppliers":       "warehouse/suppliers",
				"Products":        "warehouse/products",
				"Categories":      "warehouse/categories",
				"Warehouses":      "warehouse/warehouses",
				"Stock Movements": "warehouse/stock_movements",
				"Purchase Orders": "warehouse/purchase_orders",
			},
			"Restaurant": {
				"Menu Items":  "restaurant/menu_items",
				"Categories":  "restaurant/categories",
				"Tables":      "restaurant/tables",
				"Orders":      "restaurant/orders",
				"Order Items": "restaurant/order_items",
				"Staff":       "restaurant/staff",
			},
			"Task": {
				"Projects":    "task/projects",
				"Tasks":       "task/tasks",
				"Users":       "task/users",
				"Comments":    "task/comments",
				"Attachments": "task/attachments",
			},
		},
	}

	homeModule := &modhome.UI{
		Layout: layoutInstance,
	}

	uiSitemap := uiModule.Sitemap()
	if uiSitemap != nil {
		layoutInstance.AddSitemap(uiSitemap)
	}

	homeSitemap := homeModule.Sitemap()
	if homeSitemap != nil {
		layoutInstance.AddSitemap(homeSitemap)
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
