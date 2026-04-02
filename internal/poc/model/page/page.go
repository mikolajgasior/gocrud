package modelpage

const (
	PageStatusDraft = iota + 1
	PageStatusPublished
	PageStatusTrash
)

const (
	ElementTypeHTML = iota + 1
)

type Page struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title" crud:"req len:1,120"`
	Description string `json:"description" crud:"len:,120"`
	URI         string `json:"uri" crud:"req len:1,120" crud_regexp:"^[a-zA-Z0-9-_]+$"`
	StatusID    uint8  `json:"status_id" crud:"req val:1,3"`
	CreatedAt   int64  `json:"created_at"`
	CreatedBy   uint64 `json:"created_by"`
	ModifiedAt  int64  `json:"modified_at"`
	ModifiedBy  uint64 `json:"modified_by"`
}

type Row struct {
	ID      uint64 `json:"id"`
	PageID  uint64 `json:"page_id"`
	NumCols uint8  `json:"num_cols" crud:"req val:1,5"`
	Order   int    `json:"order"`
	Width   uint8  `json:"width" crud:"req val:1,100"`
}

type Col struct {
	ID     uint64 `json:"id"`
	RowID  uint64 `json:"row_id"`
	ColNum uint8  `json:"col_num" crud:"req val:1,"`
	Width  uint8  `json:"width" crud:"req val:1,100"`
}

type Element struct {
	ID      uint64 `json:"id"`
	ColID   uint64 `json:"col_id"`
	TypeID  uint64 `json:"type_id" crud:"req val:1,1"`
	Content string `json:"content"`
}
