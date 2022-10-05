package service

import (
	"database/sql"
	"testcode/test3/notifier"
)

type Service struct {
	Db       *sql.DB
	Notifier notifier.Notifier
}
