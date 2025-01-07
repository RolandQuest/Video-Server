package VideoDatabase

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var (
	Handle *sql.DB
	Version string
)

func Initialize(dbPath string) error {
	var err error
	if Handle, err = sql.Open("sqlite3", dbPath); err != nil {
		return err
	}
	Handle.QueryRow(`SELECT sqlite_version()`).Scan(&Version)
	return nil
}

type Tag struct {
	Id int64
	Name string
}

type Video struct {
	Id int64
	Title string
	Size int64
}

type Season struct {
	Id int64
	Name string
	Videos []Video
}

type Series struct {
	Id int64
	Name string
	Tags []Tag
	Seasons []Season
}

func SQLVersion() string {
	return Version
}

