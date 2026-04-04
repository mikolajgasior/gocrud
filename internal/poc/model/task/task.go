package modeltask

// Project - Container for related tasks
type Project struct {
	ID          uint64  `json:"id"`
	Code        string  `json:"code" crud:"req len:3,20 crud_regexp:\"^[A-Z0-9-]+$\""`
	Name        string  `json:"name" crud:"req len:2,150"`
	Description string  `json:"description" crud:"len:0,1000"`
	OwnerID     uint64  `json:"owner_id" crud:"val:1"`
	StartDate   int64   `json:"start_date" crud:"val:0"`
	EndDate     int64   `json:"end_date" crud:"val:0"`
	Status      string  `json:"status" crud:"req len:3,20"`
	Priority    string  `json:"priority" crud:"req len:3,20"`
	Budget      float64 `json:"budget" crud:"val:0"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   int64   `json:"created_at"`
	CreatedBy   uint64  `json:"created_by"`
	ModifiedAt  int64   `json:"modified_at"`
	ModifiedBy  uint64  `json:"modified_by"`
}

// Task - Individual work items
type Task struct {
	ID            uint64  `json:"id"`
	TaskNumber    string  `json:"task_number" crud:"req len:5,30 crud_regexp:\"^TASK-[0-9]+$\""`
	Title         string  `json:"title" crud:"req len:3,200"`
	Description   string  `json:"description" crud:"len:0,2000"`
	ProjectID     uint64  `json:"project_id" crud:"val:1"`
	ParentTaskID  uint64  `json:"parent_task_id" crud:"val:0"`
	AssignedTo    uint64  `json:"assigned_to" crud:"val:1"`
	ReporterID    uint64  `json:"reporter_id" crud:"val:1"`
	Status        string  `json:"status" crud:"req len:3,20"`
	Priority      string  `json:"priority" crud:"req len:3,20"`
	DueDate       int64   `json:"due_date" crud:"val:0"`
	EstimateHours float64 `json:"estimate_hours" crud:"val:0"`
	ActualHours   float64 `json:"actual_hours" crud:"val:0"`
	Progress      int     `json:"progress" crud:"val:0,100"`
	CreatedAt     int64   `json:"created_at"`
	CreatedBy     uint64  `json:"created_by"`
	ModifiedAt    int64   `json:"modified_at"`
	ModifiedBy    uint64  `json:"modified_by"`
}

// User - System users/team members
type User struct {
	ID         uint64 `json:"id"`
	Username   string `json:"username" crud:"req len:3,30 crud_regexp:\"^[a-z0-9_-]+$\""`
	Email      string `json:"email" crud:"req email"`
	FirstName  string `json:"first_name" crud:"req len:2,50"`
	LastName   string `json:"last_name" crud:"req len:2,50"`
	Phone      string `json:"phone" crud:"len:7,20"`
	Department string `json:"department" crud:"len:0,100"`
	Role       string `json:"role" crud:"req len:3,30"`
	IsActive   bool   `json:"is_active"`
	LastLogin  int64  `json:"last_login" crud:"val:0"`
	CreatedAt  int64  `json:"created_at"`
	CreatedBy  uint64 `json:"created_by"`
	ModifiedAt int64  `json:"modified_at"`
	ModifiedBy uint64 `json:"modified_by"`
}

// Comment - Task comments and discussions
type Comment struct {
	ID              uint64 `json:"id"`
	TaskID          uint64 `json:"task_id" crud:"val:1"`
	UserID          uint64 `json:"user_id" crud:"val:1"`
	Content         string `json:"content" crud:"req len:1,2000"`
	IsEdited        bool   `json:"is_edited"`
	ParentCommentID uint64 `json:"parent_comment_id" crud:"val:0"`
	CreatedAt       int64  `json:"created_at"`
	CreatedBy       uint64 `json:"created_by"`
	ModifiedAt      int64  `json:"modified_at"`
	ModifiedBy      uint64 `json:"modified_by"`
}

// Attachment - Files linked to tasks
type Attachment struct {
	ID          uint64 `json:"id"`
	TaskID      uint64 `json:"task_id" crud:"val:1"`
	FileName    string `json:"file_name" crud:"req len:3,255"`
	FilePath    string `json:"file_path" crud:"req len:5,500"`
	FileSize    int64  `json:"file_size" crud:"val:1"`
	FileType    string `json:"file_type" crud:"len:3,50"`
	UploadedBy  uint64 `json:"uploaded_by" crud:"val:1"`
	Description string `json:"description" crud:"len:0,300"`
	CreatedAt   int64  `json:"created_at"`
	CreatedBy   uint64 `json:"created_by"`
	ModifiedAt  int64  `json:"modified_at"`
	ModifiedBy  uint64 `json:"modified_by"`
}
