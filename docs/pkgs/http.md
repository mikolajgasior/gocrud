# HTTP API

The `pkg/http/api` package mounts a ready-made JSON REST API on top of a [Service](service.md) instance. A single handler covers list, read, create/update, and delete for every path registered in the service.

## Initialization

```go
import (
    "codeberg.org/mikolajgasior/gocrud/pkg/http/api"
    svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

svc := svccrud.New(paths, dbConn, gocrud.DialectPostgres)

handler := api.New(svc, api.Options{})
```

`New` takes two arguments:

| Argument | Type | Description |
|---|---|---|
| `svc` | `*service.CRUD` | An initialised service instance |
| `options` | `api.Options` | Configuration for the handler (zero value is valid) |

### Options

```go
type Options struct {
    CORS  cors.CORS
    Paths map[string]PathOptions
}
```

| Field | Type | Description |
|---|---|---|
| `CORS` | `cors.CORS` | CORS headers written on every response. Zero value emits no CORS headers. |
| `Paths` | `map[string]PathOptions` | Per-path configuration. Paths absent from the map use the zero value (all operations enabled, all filters allowed). |

### PathOptions

```go
type PathOptions struct {
    // Operation guards
    DisableCreate  bool
    DisableUpdate  bool
    DisableDelete  bool
    DisableRead    bool
    DisableList    bool

    // Filter constraints
    DisableFilters bool
    AllowedFilters []string

    // Per-operation constructor overrides
    CreateConstructor func() interface{}
    UpdateConstructor func() interface{}
    ReadConstructor   func() interface{}
    ListConstructor   func() interface{}
}
```

**Operation guards**

| Field | Type | Description |
|---|---|---|
| `DisableCreate` | `bool` | Reject `PUT /{path}/` (create) with `405` |
| `DisableUpdate` | `bool` | Reject `PUT /{path}/{id}` (update) with `405` |
| `DisableDelete` | `bool` | Reject `DELETE /{path}/{id}` with `405` |
| `DisableRead` | `bool` | Reject `GET /{path}/{id}` with `405` |
| `DisableList` | `bool` | Reject `GET /{path}/` with `405` |

**Filter constraints**

| Field | Type | Description |
|---|---|---|
| `DisableFilters` | `bool` | Ignore all `filter_val_*` / `filter_op_*` query parameters |
| `AllowedFilters` | `[]string` | Whitelist of field names that may be used as filters. Empty slice means all fields are allowed (unless `DisableFilters` is set). |

**Constructor overrides**

By default each operation uses the constructor registered for the path in the service. Setting a constructor here overrides that for the specific operation, allowing a different struct type to be used — for example a create-only struct with fewer fields that implements a custom `InsertQuery()` method.

| Field | Type | Description |
|---|---|---|
| `CreateConstructor` | `func() interface{}` | Constructor used when creating a new record (`PUT /{path}/`) |
| `UpdateConstructor` | `func() interface{}` | Constructor used when updating a record (`PUT /{path}/{id}`). The existing record is **not** pre-loaded; the URL id is stamped onto the object after JSON unmarshalling. |
| `ReadConstructor` | `func() interface{}` | Constructor used when reading a single record (`GET /{path}/{id}`). Only the fields present on the override struct are SELECTed. |
| `ListConstructor` | `func() interface{}` | Constructor used when listing records (`GET /{path}/`). Only the fields present on the override struct are SELECTed. |

**Example** — read-only path, restricted filtering, and a minimal create struct:

```go
// The "_" suffix is stripped when deriving the table name:
//   strings.Split("User_Create", "_")[0] = "User" → table "user"
type User_Create struct {
    ID    uint64
    Email string `crud:"req email"`
    Role  string `crud:"req len:3,30"`
}

handler := api.New(svc, api.Options{
    Paths: map[string]api.PathOptions{
        "users": {
            DisableDelete:     true,
            AllowedFilters:    []string{"Role", "IsActive"},
            CreateConstructor: func() interface{} { return &User_Create{} },
        },
        "audit_logs": {
            DisableCreate:  true,
            DisableUpdate:  true,
            DisableDelete:  true,
            DisableFilters: true,
        },
    },
})
```

## Mounting

`handler.Serve` is a standard `http.HandlerFunc`. Mount it under a prefix using `http.StripPrefix` so the handler only sees the path-relative URL:

```go
mux := http.NewServeMux()

subMux := http.NewServeMux()
subMux.HandleFunc("/", handler.Serve)
mux.Handle("/api/", http.StripPrefix("/api", subMux))

http.ListenAndServe(":8080", mux)
```

After stripping the mount prefix the handler expects URLs of the form `/{path}/{id}` where `path` is a key from the registry and `id` is an optional numeric record ID.

## URL scheme

```
/{path}/           — list / create
/{path}/{id}       — read / update / delete
```

`path` must exactly match a key in the paths registry (e.g. `users`, `warehouse/products`).
`id` must be a positive integer or absent.

## Password fields

Fields tagged `crud:"pass"` are **never included in responses** from the Read and List endpoints. The handler strips them before serialising the object to JSON, regardless of whether the struct field has a `json` tag or not. This applies to both the default struct and any override struct supplied via `ReadConstructor` / `ListConstructor`.

```go
type User struct {
    ID       uint64
    Email    string `crud:"req email"`
    Password string `crud:"pass"` // omitted from GET /users/ and GET /users/{id}
}
```

The CRUD layer also zeroes password fields immediately after scanning a row from the database (in both `Load` and `Get`), so hashes never reach caller code even outside the HTTP layer.

## Endpoints

### List — `GET /{path}/`

Returns a paginated, filtered list of records.

**Query parameters:**

| Parameter | Type | Default | Description |
|---|---|---|---|
| `limit` | int | `10` | Maximum number of records to return |
| `offset` | int | `0` | Number of records to skip |
| `order` | string | — | Field name to sort by |
| `order_direction` | string | — | `asc` or `desc` |
| `filter_val_{Field}` | string | — | Filter value for the named field |
| `filter_op_{Field}` | string | `eq` | Filter operator for the named field (see [operators](service.md#list)) |

**Example request:**
```
GET /users/?limit=20&offset=0&order=LastName&order_direction=asc&filter_val_Role=admin&filter_op_Role=eq
```

**Response `200 OK`:**
```json
{
  "ok": true,
  "code": "SUCCESS",
  "data": [ { "id": 1, "username": "alice", ... }, ... ]
}
```

---

### Read — `GET /{path}/{id}`

Returns a single record by ID.

**Response `200 OK`:**
```json
{
  "ok": true,
  "code": "SUCCESS",
  "data": { "id": 42, "username": "alice", ... }
}
```

**Response `404 Not Found`:**
```json
{ "ok": true, "code": "NOT_FOUND" }
```

---

### Create — `PUT /{path}/`

Creates a new record. The request body must be a JSON object whose keys match the struct's `json` tags.

**Request body:**
```json
{ "username": "alice", "email": "alice@example.com", "role": "admin" }
```

**Response `201 Created`:**
```json
{
  "ok": true,
  "code": "CREATED",
  "data": { "id": 42 }
}
```

**Response `400 Bad Request`** (validation failure):
```json
{
  "ok": true,
  "code": "VALIDATION_FAILED",
  "message": "validation failed with violations: ..."
}
```

---

### Update — `PUT /{path}/{id}`

Updates an existing record. The request body replaces the record's fields. The record must exist; if not, `404` is returned.

**Response `200 OK`:**
```json
{ "ok": true, "code": "SUCCESS", "data": { "id": 42 } }
```

---

### Delete — `DELETE /{path}/{id}`

Deletes a record by ID. Cascade rules defined on the struct are applied automatically.

**Response `200 OK`:**
```json
{ "ok": true, "code": "SUCCESS" }
```

**Response `404 Not Found`:**
```json
{ "ok": true, "code": "NOT_FOUND" }
```

---

## Response format

Every response is a JSON object with the following shape:

```json
{
  "ok":      true,
  "code":    "SUCCESS",
  "message": "optional human-readable detail",
  "data":    null
}
```

| Field | Type | Description |
|---|---|---|
| `ok` | bool | Always `true` — use `code` to distinguish outcomes |
| `code` | string | Machine-readable result code (see below) |
| `message` | string? | Present only on errors; contains a human-readable detail |
| `data` | any? | Present on successful read / list / create / update |

**Response codes:**

| Code | HTTP status | Meaning |
|---|---|---|
| `SUCCESS` | 200 | Operation completed successfully |
| `CREATED` | 201 | Record created |
| `NOT_FOUND` | 404 | No record with the given ID |
| `VALIDATION_FAILED` | 400 | Request body failed struct validation |
| `BAD_REQUEST` | 400 | Malformed URL or unparseable request body |
| `URL_PATH_ID` | 400 | ID segment in the URL is not a valid number |
| `NOT_ALLOWED` | 405 | Operation is disabled for this path via `PathOptions` |
| `SERVICE_ERROR` | 500 | Internal error from the service layer |

## CORS

Populate `Options.CORS` to emit CORS headers on every response:

```go
import (
    "codeberg.org/mikolajgasior/gocrud/pkg/http/api"
    "codeberg.org/mikolajgasior/gocrud/pkg/http/cors"
)

handler := api.New(svc, api.Options{
    CORS: cors.CORS{
        AllowOrigin:  "https://app.example.com",
        AllowHeaders: "Content-Type, Authorization",
        AllowMethods: "GET, PUT, DELETE",
        MaxAge:       3600,
    },
})
```

Leave `CORS` as its zero value (`cors.CORS{}`) to emit no CORS headers.

## Complete example

```go
package main

import (
    "net/http"

    "codeberg.org/mikolajgasior/gocrud"
    "codeberg.org/mikolajgasior/gocrud/pkg/http/api"
    svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
    _ "github.com/lib/pq"
)

type User struct {
    ID    uint64 `json:"id"`
    Name  string `json:"name"  crud:"req len:2,100"`
    Email string `json:"email" crud:"req email"`
}

func main() {
    db, _ := sql.Open("postgres", "host=localhost user=app password=secret dbname=app sslmode=disable")

    paths := map[string]func() interface{}{
        "users": func() interface{} { return &User{} },
    }

    svc := svccrud.New(paths, db, gocrud.DialectPostgres)
    handler := api.New(svc, api.Options{})

    mux := http.NewServeMux()
    subMux := http.NewServeMux()
    subMux.HandleFunc("/", handler.Serve)
    mux.Handle("/api/", http.StripPrefix("/api", subMux))

    http.ListenAndServe(":8080", mux)
}
```

With this setup:

| Request | Action |
|---|---|
| `GET /api/users/` | List users |
| `GET /api/users/1` | Read user 1 |
| `PUT /api/users/` | Create user |
| `PUT /api/users/1` | Update user 1 |
| `DELETE /api/users/1` | Delete user 1 |
