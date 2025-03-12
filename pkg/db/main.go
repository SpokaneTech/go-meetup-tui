package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/SpokaneTech/go-meetup-tui/pkg/drivers"
)

type DB interface {
	Open() error
	Close() error

	Query(query string, args ...interface{}) ([]string, []interface{}, error)
	Exec(query string, args ...interface{}) error
}

func New(dsn string) (DB, error) {
	driver, _, found := strings.Cut(dsn, ":")
	if !found {
		return nil, errors.New("invalid dsn")
	}

	switch driver {
	case "sqlite":
		return drivers.Sqlite(dsn), nil

	default:
		return nil, fmt.Errorf("unknown driver %s", driver)
	}
}
