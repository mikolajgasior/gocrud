package crud

import (
	"database/sql"
	"sync"

	sqlbuilder "codeberg.org/mikolajgasior/gocrud/pkg/sql/builder"
)

type CRUD struct {
	db              *sql.DB
	buildersMu      sync.RWMutex
	builders        map[string]*sqlbuilder.Builder
	tagName         string
	tableNamePrefix string
	flags           uint64
}

func New(db *sql.DB, options Options) *CRUD {
	crud := &CRUD{
		db:       db,
		builders: make(map[string]*sqlbuilder.Builder),
	}

	crud.tagName = "crud"
	if options.TagName != "" {
		crud.tagName = options.TagName
	}

	if options.TableNamePrefix != "" {
		crud.tableNamePrefix = options.TableNamePrefix
	}

	crud.flags = options.Flags

	return crud
}

func (c *CRUD) SetFlag(flag uint64) {
	if c.flags&flag == 0 {
		c.flags |= flag
	}
}
