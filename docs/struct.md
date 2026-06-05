# Struct

To utilize the library's automatic CRUD capabilities, your Go structs must adhere to a specific naming convention and tagging strategy. These definitions allow the library to map your data structures to database records and RESTful API endpoints while enforcing validation rules automatically.

## Basic requirements

### 1. The ID Field

Every struct intended for CRUD operations must include an ID field which must be uint64. This serves as the primary key for database operations.

### 2. Supported Data Types

At this stage, the library supports the following basic data types for struct fields:

* Integers: int, int8, int16, int32, int64
* Unsigned Integers: uint, uint8, uint16, uint32, uint64
* Floats: float32, float64
* Strings: string
* Booleans: bool

Complex types (such as slices, maps, or nested structs) are not currently supported for direct mapping in this context.

## Tagging

To configure field behavior, you will use two primary struct tags: json and crud.

### The `json` Tag

Use the `json` tag to define the field name as it appears in your RESTful API payloads. This ensures your API follows standard naming conventions (e.g., snake_case) while your Go code uses idiomatic naming (PascalCase).

* **Format**: `json:"api_field_name"`
* **Example**: `json:"first_name"` maps the Go field FirstName to first_name in JSON.

### The `crud` Tag

The `crud` tag is the core mechanism for validation. It defines rules that are enforced in two specific scenarios:

* **Saving Data**: Before a record is saved to the database via the `Save` function.
* **Filtering Data**: When using filters in the `Get` method to ensure query parameters are valid.

These tags leverage the underlying validation logic provided by the [struct-validator](https://github.com/mikolajgasior/struct-validator) library.

#### Common Validation Rules

You can combine multiple rules within a single crud tag string.

* `req`: Marks the field as required.
* `len:min,max`: Validates string length.
* `val:min,max`: Validates numeric range.
* `email`: Validates that the string is a properly formatted email address.
* `crud_regexp`:"pattern": Validates the field against a regular expression.
* `pass`: Marks a `string` field as a password. On **Save** the plain-text value is automatically bcrypt-hashed before it is written to the database. On **Load** and **Get** the field is zeroed after the row is scanned so that bcrypt hashes never appear in memory after a read. The HTTP API (`pkg/http/api`) goes further and omits the field entirely from JSON responses.

At least one of `min`, `max` must be specified.

**Note on Special Characters:** When using regular expressions inside struct tags, remember to escape double quotes. For example: `crud_regexp:"^[a-z]+$"`.

## Complete Example

Below is a comprehensive example of a `User` struct that demonstrates all the requirements and tagging strategies discussed.

```go
type User struct {
	// Required ID field (uint64)
	ID         uint64 `json:"id"`

	// Username: Required, length 3-30, must match regex
	Username   string `json:"username" crud:"req len:3,30 crud_regexp:\"^[a-z0-9_-]+$\""`

	// Email: Required, must be valid email format
	Email      string `json:"email" crud:"req email"`

	// First Name: Required, length 2-50
	FirstName  string `json:"first_name" crud:"req len:2,50"`

	// Last Name: Required, length 2-50
	LastName   string `json:"last_name" crud:"req len:2,50"`

	// Phone: Optional, length 7-20
	Phone      string `json:"phone" crud:"len:7,20"`

	// Department: Optional, max length 100 (empty string allowed)
	Department string `json:"department" crud:"len:0,100"`

	// Role: Required, length 3-30
	Role       string `json:"role" crud:"req len:3,30"`

	// IsActive: Boolean flag (no specific validation rules applied here)
	IsActive   bool   `json:"is_active"`

	// LastLogin: Must be 0 or a positive timestamp
	LastLogin  int64  `json:"last_login" crud:"val:0"`

	// Timestamps and Audit Fields
	CreatedAt  int64  `json:"created_at"`
	CreatedBy  uint64 `json:"created_by"`
	ModifiedAt int64  `json:"modified_at"`
	ModifiedBy uint64 `json:"modified_by"`
}
```