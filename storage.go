package main

import (
	"database/sql"
	"fmt"
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
		fmt.Print(err)
		panic(err)
	}

	statement := `
	CREATE TABLE IF NOT EXISTS builds(
		id INTEGER NOT NULL PRIMARY KEY,
		name STRING NOT NULL UNIQUE,
		path STRING NOT NULL,
		branch STRING);`

	_, err = db.Exec(statement)
	if err != nil {
		fmt.Printf("%q: %s", err, statement)
		panic(err)
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

	_, err := db.Exec(statement, info.Name, info.AbsolutePath, info.Branch)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteBuildInfo(db *sql.DB, name string) {
	statement := `
	DELETE FROM builds
	WHERE name = ?`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	_, err := db.Exec(statement, name)
	if err != nil {
		fmt.Println(err)
	}
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

		fmt.Println(err)
	}

	return &build
}

func getBuilds(db *sql.DB) []*buildInfo {
	statement := `
	SELECT * FROM builds`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err)
	}

	builds := make([]*buildInfo, 0)

	for rows.Next() {
		build := buildInfo{}
		if err := rows.Scan(&build.ID,
			&build.Name,
			&build.AbsolutePath,
			&build.Branch); err != nil {

			fmt.Println(err)
		}
		builds = append(builds, &build)
	}

	return builds
}
