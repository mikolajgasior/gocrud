package main

import (
	"context"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/app"
	modelrestaurant "codeberg.org/mikolajgasior/gocrud/internal/poc/model/restaurant"
	modeltask "codeberg.org/mikolajgasior/gocrud/internal/poc/model/task"
	modelwarehouse "codeberg.org/mikolajgasior/gocrud/internal/poc/model/warehouse"
	"codeberg.org/mikolajgasior/gocrud/internal/poc/module"
	modapi "codeberg.org/mikolajgasior/gocrud/internal/poc/module/api"
)

func main() {
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

	var appObj = app.App{
		Modules: map[string]module.Module{
			"10_api": apiModule,
		},
	}
	appObj.Run(context.Background())
}
