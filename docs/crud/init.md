# Initialization

To begin using the CRUD library, you must initialize a new instance by passing your database connection and optional configuration settings. This creates the controller object that manages all subsequent database operations.

## Basic Initialization

Create a new CRUD instance using the `New` function:

```go
crudInstance := structcrud.New(dbConn, structcrud.Options{})
```

## Configuration Options

The `Options` struct allows you to customize the behavior of the CRUD instance during initialization:

```go
type Options struct {
    TableNamePrefix string
    TagName         string
    Flags           uint64
}
```

* **TableNamePrefix (string)**: Allows you to define a prefix for all database table names managed by this instance. This is useful for organizing tables in shared databases or separating environments (e.g., dev_, prod_). If left empty, the library uses the default table name derived from the struct.
* **TagName (string)**: Specifies the name of the struct tag used for validation and mapping rules. Default: `crud`. If you set this to `mytag`, the library will look for tags like `mytag:"unique"` instead of the default `crud:"unique"`.

* **Flags (uint64)**: Enables specific behaviors via bitwise flags. Currently, one flag is supported: **GetCountOnUniq (const GetCountOnUniq = 1)**: When this flag is set, the library performs an additional COUNT(*) query on unique fields before attempting an insert or update. This ensures that the unique value does not already exist in the database, providing an extra layer of validation beyond standard database constraints.

## Example Usage

Below is a complete example initializing a CRUD instance with a table prefix and enabling the uniqueness count check:

```go
import (
    structcrud "codeberg.org/mikolajgasior/gocrud"
)

// Create a database connection
// for example: dbConn := db.Open("mysql", "root:password@tcp(localhost:3306)/mydb")

// Define options
opts := structcrud.Options{
    TableNamePrefix: "app_",
    TagName:         "crud",
    Flags:           structcrud.GetCountOnUniq,
}

// Initialize the instance
crudInstance := structcrud.New(dbConn, opts)

// Now crudInstance is ready to perform operations
// Tables will be prefixed with "app_"
// Unique checks will run a COUNT(*) query before insertion
```

