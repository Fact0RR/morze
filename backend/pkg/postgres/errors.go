package postgres

const (
	InvalidConnectionErrMsg  = "invalid postgres connection URI"
	UnablePoolErrMsg         = "unable to create connection pool"
	CannotTalkPostgresErrMsg = "cannot talk to postgres, last attempt failed"
	CannotGetPGPoolErrMsg    = "cannot get PG pool"
	ErrorFmtTemplate         = "%s: %w"
)
