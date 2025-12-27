package crud

import (
	"database/sql"
	"regexp"

	sqlbuilder "github.com/keenbytes/pgsql-builder"
)

const (
	IDField = "ID"
)

var (
	tagWithValRegexp = regexp.MustCompile(`[a-zA-Z0-9_]+:[a-zA-Z0-9_-]+`)
)

type CRUD struct {
	db              *sql.DB
	builders        map[string]*sqlbuilder.Builder
	tagName         string
	tableNamePrefix string
	flags           int64
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

func (c *CRUD) SetFlag(flag int64) {
	if c.flags&flag == 0 {
		c.flags |= flag
	}
}
