package app

const (
	LogInvalidPort   = "invalid port"
	LogCannotOpenDB  = "failed to connect to db"
	LogCannotCloseDB = "failed to close db connection"
	LogHTTPListen    = "failed to start http server"
)

const (
	ExitInvalidPort = iota
	ExitCannotOpenDB
	ExitCannotCloseDB
	ExitHTTPListen
	ExitInitError
	ExitCreateTableError
)
