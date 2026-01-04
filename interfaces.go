package crud

import sqlbuilder "github.com/keenbytes/pgsql-builder"

type deleteByIDer interface {
	DeleteByID() (string, error)
}

type deleteReturningIDer interface {
	DeleteReturningID() (string, error)
}

type selecter interface {
	Select(order []string, limit int, offset int, filters *sqlbuilder.Filters) (string, error)
}

type selectCounter interface {
	SelectCount(filters *sqlbuilder.Filters) (string, error)
}

type selectByIDer interface {
	SelectByID() (string, error)
}

type updateByIDer interface {
	UpdateByID() (string, error)
}

type insertOnConflictUpdateer interface {
	InsertOnConflictUpdate() (string, error)
}

type inserter interface {
	Insert() (string, error)
}

type updateer interface {
	Update(values map[string]interface{}, filters *sqlbuilder.Filters) (string, error)
}

type createTableer interface {
	CreateTable() (string, error)
}

type dropTableer interface {
	DropTable() (string, error)
}
