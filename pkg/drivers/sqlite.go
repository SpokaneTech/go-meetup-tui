package drivers

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func Sqlite(dsn string) *sqlitedb {
	return &sqlitedb{path: dsn}
}

type sqlitedb struct {
	path string
	db   *sql.DB
}

func (s *sqlitedb) Open() error {
	path, _ := strings.CutPrefix(s.path, "sqlite:")

	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	_, err := os.Stat(filepath.Dir(path))
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dir, os.ModePerm)
		} else {
			return err
		}
	}

	db, err := sql.Open("sqlite", filepath.Join(dir, filename))
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *sqlitedb) Close() error {
	return s.db.Close()
}

func (s *sqlitedb) Query(query string, args ...interface{}) ([]string, []interface{}, error) {
	r, err := s.db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}

	return scanRows(r)
}

func (s *sqlitedb) Exec(query string, args ...interface{}) error {
	_, err := s.db.Exec(query, args...)

	return err
}
