package querycontainer

type DropTableBuilder interface {
	BuildDropTableQuery(tableNamePrefix string) (string, error)
}

type CreateTableBuilder interface {
	BuildCreateTableQuery(tableNamePrefix string) (string, error)
}

type DeleteByIDBuilder interface {
	BuildDeleteByIDQuery(tableNamePrefix string) (string, error)
}

type DeletePrefixBuilder interface {
	BuildDeletePrefixQuery(tableNamePrefix string) (string, error)
}

type UpdateByIDBuilder interface {
	BuildUpdateByIDQuery(tableNamePrefix string) (string, error)
}

type UpdatePrefixBuilder interface {
	BuildUpdatePrefixQuery(tableNamePrefix string) (string, error)
}

type InsertBuilder interface {
	BuildInsertQuery(tableNamePrefix string) (string, error)
}

type InsertOnConflictUpdateBuilder interface {
	BuildInsertOnConflictUpdateQuery(tableNamePrefix string) (string, error)
}

type SelectByIDBuilder interface {
	BuildSelectByIDQuery(tableNamePrefix string) (string, error)
}

type SelectPrefixBuilder interface {
	BuildSelectPrefixQuery(tableNamePrefix string) (string, error)
}

type SelectCountPrefixBuilder interface {
	BuildSelectCountPrefixQuery(tableNamePrefix string) (string, error)
}
