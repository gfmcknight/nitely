package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type buildInfo struct {
	ID           int
	Name         string
	AbsolutePath string
	Branch       string
}

type serviceInfo struct {
	Name         string
	AbsolutePath string
}

func getStorageBase() string {
	return os.Getenv("NitelyPath")
}

func openAndCreateStorage() *sql.DB {
	db, err := sql.Open("sqlite3", filepath.Join(getStorageBase()+"/build-info.db"))
	if err != nil {
		log.Fatal(err)
	}

	statement := `
	CREATE TABLE IF NOT EXISTS builds(
		id INTEGER NOT NULL PRIMARY KEY,
		name STRING NOT NULL UNIQUE,
		path STRING NOT NULL,
		branch STRING);`

	_, err = db.Exec(statement)
	if err != nil {
		log.Printf("%q: %s", err, statement)
		return nil
	}

	return db
}

func insertBuildInfo(db *sql.DB, info buildInfo) {
	statement := `
	INSERT INTO builds(name, path, branch)
	VALUES(?, ?, ?)`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	db.Exec(statement, info.Name, info.AbsolutePath, info.Branch)
}

func deleteBuildInfo(db *sql.DB, name string) {
	statement := `
	DELETE FROM builds
	WHERE name = ?`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	db.Exec(statement, name)
}

func getBuildInfo(db *sql.DB, name string) *buildInfo {
	statement := `
	SELECT * FROM builds b
	WHERE b.name == ?`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	build := buildInfo{}
	results := db.QueryRow(statement, name)

	if err := results.Scan(&build.ID,
		&build.Name,
		&build.AbsolutePath,
		&build.Branch); err != nil {

		log.Fatal(err)
	}

	return &build
}
