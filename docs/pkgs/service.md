# Service

The `pkg/service` package provides a higher-level wrapper around the core CRUD object. Where the core layer works directly with struct instances, the service layer is driven by a **path registry** — a map of string keys to constructor functions — and accepts filter values as plain strings, making it a natural fit for wiring up to HTTP handlers.

## Initialization

```go
import (
    "codeberg.org/mikolajgasior/gocrud"
    svccrud    "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

svc := svccrud.New(
    map[string]func() interface{}{
        "users":    func() interface{} { return &User{} },
        "products": func() interface{} { return &Product{} },
    },
    dbConn,
    gocrud.DialectPostgres, // or gocrud.DialectSQLite
)
```

`New` takes three arguments:

* **paths** `map[string]func() interface{}` — a registry mapping path keys to constructor functions. Each function must return a pointer to a new zero-value struct.
* **dbConn** `*sql.DB` — an open database connection.
* **dialect** `string` — the SQL dialect; must be `gocrud.DialectPostgres` or `gocrud.DialectSQLite`. Passing an empty or unrecognised value causes a panic at startup.

### Creating tables

```go
err := svc.CreateTables(ctx)
```

Iterates over all registered paths and calls `CREATE TABLE IF NOT EXISTS` for each struct. Useful during application startup when `CREATE_TABLES=true` is set.

## Methods

### Save

```go
err := svc.Save(ctx, obj, now, userID)
```

Saves `obj` to the database. If `obj.ID == 0` an `INSERT` is performed; otherwise an upsert. `now` (Unix timestamp) and `userID` are written to the audit fields (`CreatedAt`, `CreatedBy`, `ModifiedAt`, `ModifiedBy`) when those fields are present on the struct.

**Returns** `*ModelValidationError` when struct validation fails, otherwise a generic error.

---

### SaveFromForm

```go
err := svc.SaveFromForm(ctx, obj, urlValues, "prefix_", now, userID)
```

Populates `obj` from `url.Values` (e.g. an HTTP form submission) before saving. Each form key has `namePrefix` stripped from the front, and the remaining string is matched against struct field names. Type conversion from string to the field's native type is applied automatically.

**Returns** `*ModelValidationError` if any value cannot be converted to the target field type, or if struct validation fails after population.

---

### Read

```go
obj, err := svc.Read(ctx, "users", id)
```

Loads a single record by `id`. Returns `NotFoundError` if no record with that ID exists.

---

### Delete

```go
err := svc.Delete(ctx, "users", id)
```

Reads the record by `id` first (returning `NotFoundError` if absent), then deletes it. Cascade rules defined on the struct are applied automatically.

---

### List

```go
objs, err := svc.List(
    ctx,
    "users",
    limit, offset,
    "Age", "asc",
    map[string]string{"Status": "1"},
    map[string]string{"Status": "eq"},
    nil, // optional row transform func
)
```

Returns a paginated, filtered, ordered list of objects for the given path. All filter values are passed as strings and converted to the correct field types automatically.

**Parameters:**

| Parameter | Type | Description |
|---|---|---|
| `path` | `string` | Registry key |
| `limit` | `int` | SQL `LIMIT` |
| `offset` | `int` | SQL `OFFSET` |
| `order` | `string` | Field name to sort by (empty = no ordering) |
| `orderDirection` | `string` | `"asc"` or `"desc"` |
| `filterVals` | `map[string]string` | Field name → value |
| `filterOps` | `map[string]string` | Field name → operator string (see below) |
| `rowFunc` | `func(interface{}) interface{}` | Optional per-row transform; `nil` returns raw structs |

**Filter operators:**

| String | Meaning |
|---|---|
| `eq` | `=` (default when omitted) |
| `ne` | `!=` |
| `lt` | `<` |
| `le` | `<=` |
| `gt` | `>` |
| `ge` | `>=` |
| `like` | `LIKE` |
| `match` | `~` (regex, PostgreSQL only) |
| `bit` | bitwise AND `> 0` |

**Returns** `*FilterValidationError` on an unknown operator or a filter field that fails struct validation.

---

### Num

```go
count, err := svc.Num(ctx, "users",
    map[string]string{"Status": "1"},
    map[string]string{"Status": "eq"},
)
```

Returns the count of records matching the given filters. Accepts the same `filterVals`/`filterOps` maps as `List`. Useful for pagination — combine with `List`'s `limit`/`offset` to get a total alongside a page of results.

---

### Helper methods

```go
obj  := svc.New("users")   // returns a new zero-value *User, or nil if path unknown
id   := svc.ID(obj)        // returns obj.ID as uint64
```

## Error types

| Type | Sentinel / constructor | When returned |
|---|---|---|
| `*ModelValidationError` | — | `Save` / `SaveFromForm`: struct field validation failed |
| `*FilterValidationError` | — | `List` / `Num`: unknown filter operator or filter validation failed |
| `NotFoundError` | `errors.Is(err, NotFoundError)` | `Read` / `Delete`: no record with the given ID |
| `InvalidPathError` | `errors.Is(err, InvalidPathError)` | Any method: path key not in the registry |

Both `ModelValidationError` and `FilterValidationError` carry a `Violations map[string]uint64` field containing the per-field failure codes from the [struct-validator](https://github.com/mikolajgasior/struct-validator) library.

## Example

The `cmd/poc` application wires the service into HTTP modules:

```go
svc := svccrud.New(map[string]func() interface{}{
    "warehouse/products":  func() interface{} { return &Product{} },
    "warehouse/suppliers": func() interface{} { return &Supplier{} },
}, dbConn, gocrud.DialectPostgres)

// In an HTTP handler:
objs, err := svc.List(ctx, "warehouse/products", 20, 0, "Name", "asc",
    map[string]string{"CategoryID": r.URL.Query().Get("category_id")},
    map[string]string{"CategoryID": "eq"},
    nil,
)
```
