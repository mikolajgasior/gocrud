package modelrestaurant

// MenuItem - Individual dishes on the menu
type MenuItem struct {
	ID           uint64  `json:"id"`
	Code         string  `json:"code" crud:"req len:3,20 crud_regexp:\"^[A-Z0-9-]+$\""`
	Name         string  `json:"name" crud:"req len:2,100"`
	Description  string  `json:"description" crud:"len:0,300"`
	CategoryID   uint64  `json:"category_id" crud:"val:1"`
	Price        float64 `json:"price" crud:"val:0"`
	Cost         float64 `json:"cost" crud:"val:0"`
	PrepTime     int     `json:"prep_time" crud:"val:1"`
	IsAvailable  bool    `json:"is_available"`
	IsVegetarian bool    `json:"is_vegetarian"`
	CreatedAt    int64   `json:"created_at"`
	CreatedBy    uint64  `json:"created_by"`
	ModifiedAt   int64   `json:"modified_at"`
	ModifiedBy   uint64  `json:"modified_by"`
}

// Category - Menu sections (Appetizers, Main, Desserts, etc.)
type Category struct {
	ID           uint64 `json:"id"`
	Code         string `json:"code" crud:"req len:2,15 crud_regexp:\"^[A-Z0-9-]+$\""`
	Name         string `json:"name" crud:"req len:2,50"`
	DisplayOrder int    `json:"display_order" crud:"val:0"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    int64  `json:"created_at"`
	CreatedBy    uint64 `json:"created_by"`
	ModifiedAt   int64  `json:"modified_at"`
	ModifiedBy   uint64 `json:"modified_by"`
}

// Table - Dining tables in the restaurant
type Table struct {
	ID         uint64 `json:"id"`
	Number     string `json:"number" crud:"req len:1,10 crud_regexp:\"^[A-Z0-9-]+$\""`
	Capacity   int    `json:"capacity" crud:"val:1,20"`
	Section    string `json:"section" crud:"len:0,50"`
	IsOutdoor  bool   `json:"is_outdoor"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  int64  `json:"created_at"`
	CreatedBy  uint64 `json:"created_by"`
	ModifiedAt int64  `json:"modified_at"`
	ModifiedBy uint64 `json:"modified_by"`
}

// Order - Customer orders
type Order struct {
	ID           uint64  `json:"id"`
	OrderNumber  string  `json:"order_number" crud:"req len:8,20 crud_regexp:\"^ORD-[0-9]+$\""`
	TableID      uint64  `json:"table_id" crud:"val:1"`
	CustomerName string  `json:"customer_name" crud:"len:0,100"`
	OrderType    string  `json:"order_type" crud:"req len:3,20"`
	Status       string  `json:"status" crud:"req len:3,20"`
	Subtotal     float64 `json:"subtotal" crud:"val:0"`
	Tax          float64 `json:"tax" crud:"val:0"`
	Total        float64 `json:"total" crud:"val:0"`
	Notes        string  `json:"notes" crud:"len:0,300"`
	CreatedAt    int64   `json:"created_at"`
	CreatedBy    uint64  `json:"created_by"`
	ModifiedAt   int64   `json:"modified_at"`
	ModifiedBy   uint64  `json:"modified_by"`
}

// OrderItem - Individual items within an order
type OrderItem struct {
	ID                  uint64  `json:"id"`
	OrderID             uint64  `json:"order_id" crud:"val:1"`
	MenuItemID          uint64  `json:"menu_item_id" crud:"val:1"`
	Quantity            int     `json:"quantity" crud:"val:1"`
	UnitPrice           float64 `json:"unit_price" crud:"val:0"`
	SpecialInstructions string  `json:"special_instructions" crud:"len:0,200"`
	IsRefunded          bool    `json:"is_refunded"`
	CreatedAt           int64   `json:"created_at"`
	CreatedBy           uint64  `json:"created_by"`
	ModifiedAt          int64   `json:"modified_at"`
	ModifiedBy          uint64  `json:"modified_by"`
}

// Staff - Restaurant employees
type Staff struct {
	ID         uint64 `json:"id"`
	EmployeeID string `json:"employee_id" crud:"req len:5,20 crud_regexp:\"^[A-Z0-9-]+$\""`
	FirstName  string `json:"first_name" crud:"req len:2,50"`
	LastName   string `json:"last_name" crud:"req len:2,50"`
	Email      string `json:"email" crud:"req email"`
	Phone      string `json:"phone" crud:"len:7,20"`
	Role       string `json:"role" crud:"req len:3,30"`
	ShiftStart string `json:"shift_start" crud:"len:0,10"`
	ShiftEnd   string `json:"shift_end" crud:"len:0,10"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  int64  `json:"created_at"`
	CreatedBy  uint64 `json:"created_by"`
	ModifiedAt int64  `json:"modified_at"`
	ModifiedBy uint64 `json:"modified_by"`
}
