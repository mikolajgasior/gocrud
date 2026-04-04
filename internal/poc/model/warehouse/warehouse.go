package modelwarehouse

// Supplier - Companies that provide products
type Supplier struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name" crud:"req len:2,100"`
	ContactName string `json:"contact_name" crud:"req len:2,80"`
	Email       string `json:"email" crud:"req email"`
	Phone       string `json:"phone" crud:"req len:7,20"`
	Address     string `json:"address" crud:"len:0,250"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   int64  `json:"created_at"`
	CreatedBy   uint64 `json:"created_by"`
	ModifiedAt  int64  `json:"modified_at"`
	ModifiedBy  uint64 `json:"modified_by"`
}

// Product - Items stored in warehouse
type Product struct {
	ID           uint64  `json:"id"`
	SKU          string  `json:"sku" crud:"req len:5,50 crud_regexp:\"^[A-Z0-9-]+$\""`
	Name         string  `json:"name" crud:"req len:2,150"`
	Description  string  `json:"description" crud:"len:0,500"`
	SupplierID   uint64  `json:"supplier_id" crud:"val:1"`
	CategoryID   uint64  `json:"category_id" crud:"val:1"`
	UnitPrice    float64 `json:"unit_price" crud:"val:0"`
	ReorderLevel int     `json:"reorder_level" crud:"val:0"`
	CurrentStock int     `json:"current_stock" crud:"val:0"`
	IsActive     bool    `json:"is_active"`
	CreatedAt    int64   `json:"created_at"`
	CreatedBy    uint64  `json:"created_by"`
	ModifiedAt   int64   `json:"modified_at"`
	ModifiedBy   uint64  `json:"modified_by"`
}

// Category - Product classifications
type Category struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code" crud:"req len:2,20 crud_regexp:\"^[A-Z0-9-]+$\""`
	Name        string `json:"name" crud:"req len:2,100"`
	Description string `json:"description" crud:"len:0,300"`
	ParentID    uint64 `json:"parent_id" crud:"val:0"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   int64  `json:"created_at"`
	CreatedBy   uint64 `json:"created_by"`
	ModifiedAt  int64  `json:"modified_at"`
	ModifiedBy  uint64 `json:"modified_by"`
}

// Warehouse - Storage locations
type Warehouse struct {
	ID         uint64 `json:"id"`
	Code       string `json:"code" crud:"req len:3,20 crud_regexp:\"^[A-Z0-9-]+$\""`
	Name       string `json:"name" crud:"req len:2,100"`
	Address    string `json:"address" crud:"req len:5,250"`
	Capacity   int    `json:"capacity" crud:"val:1"`
	ManagerID  uint64 `json:"manager_id" crud:"val:1"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  int64  `json:"created_at"`
	CreatedBy  uint64 `json:"created_by"`
	ModifiedAt int64  `json:"modified_at"`
	ModifiedBy uint64 `json:"modified_by"`
}

// StockMovement - Inventory transactions
type StockMovement struct {
	ID           uint64 `json:"id"`
	ProductID    uint64 `json:"product_id" crud:"val:1"`
	WarehouseID  uint64 `json:"warehouse_id" crud:"val:1"`
	MovementType string `json:"movement_type" crud:"req len:3,20"`
	Quantity     int    `json:"quantity" crud:"val:1"`
	Reference    string `json:"reference" crud:"len:0,100"`
	Notes        string `json:"notes" crud:"len:0,500"`
	CreatedAt    int64  `json:"created_at"`
	CreatedBy    uint64 `json:"created_by"`
	ModifiedAt   int64  `json:"modified_at"`
	ModifiedBy   uint64 `json:"modified_by"`
}

// PurchaseOrder - Orders to suppliers
type PurchaseOrder struct {
	ID           uint64  `json:"id"`
	OrderNumber  string  `json:"order_number" crud:"req len:8,30 crud_regexp:\"^PO-[0-9]+$\""`
	SupplierID   uint64  `json:"supplier_id" crud:"val:1"`
	WarehouseID  uint64  `json:"warehouse_id" crud:"val:1"`
	OrderDate    int64   `json:"order_date" crud:"val:0"`
	ExpectedDate int64   `json:"expected_date" crud:"val:0"`
	Status       string  `json:"status" crud:"req len:3,20"`
	TotalAmount  float64 `json:"total_amount" crud:"val:0"`
	Notes        string  `json:"notes" crud:"len:0,500"`
	CreatedAt    int64   `json:"created_at"`
	CreatedBy    uint64  `json:"created_by"`
	ModifiedAt   int64   `json:"modified_at"`
	ModifiedBy   uint64  `json:"modified_by"`
}
