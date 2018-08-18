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
	Remote       string
}

type serviceInfo struct {
	ID           int
	Name         string
	AbsolutePath string
	Args         string
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
		branch STRING,
		remote STRING);`
	_, err = db.Exec(statement)
	if err != nil {
		fmt.Printf("%q: %s", err, statement)
		panic(err)
	}

	statement = `
	CREATE TABLE IF NOT EXISTS services(
		id INTEGER NOT NULL PRIMARY KEY,
		name STRING NOT NULL UNIQUE,
		path STRING NOT NULL,
		args STRING);`
	_, err = db.Exec(statement)
	if err != nil {
		fmt.Printf("%q: %s", err, statement)
		panic(err)
	}

	statement = `
	CREATE TABLE IF NOT EXISTS properties(
		id INTEGER NOT NULL PRIMARY KEY,
		name STRING NOT NULL UNIQUE,
		value STRING);`
	_, err = db.Exec(statement)
	if err != nil {
		fmt.Printf("%q: %s", err, statement)
		panic(err)
	}

	return db
}

func getProperty(db *sql.DB, name string) *string {
	statement := `SELECT value FROM properties WHERE name = ?`
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	result := db.QueryRow(statement, name)

	var value string
	err := result.Scan(&value)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		fmt.Println(err)
	}

	return &value
}

func setProperty(db *sql.DB, name, value string) {
	insert := `
	INSERT INTO properties VALUES(
		(SELECT id FROM properties WHERE name = ?),
		?, ?);`

	replace := `
	UPDATE properties SET value = ? WHERE name = ?;`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	var err error
	if getProperty(db, name) == nil {
		_, err = db.Exec(insert, name, name, value)
	} else {
		_, err = db.Exec(replace, value, name)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func insertBuildInfo(db *sql.DB, info buildInfo) {
	statement := `
	INSERT INTO builds(name, path, branch, remote)
	VALUES(?, ?, ?, ?)`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	_, err := db.Exec(statement, info.Name, info.AbsolutePath, info.Branch, info.Remote)
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
		&build.Branch,
		&build.Remote); err != nil {

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
			&build.Branch,
			&build.Remote); err != nil {

			fmt.Println(err)
		}
		builds = append(builds, &build)
	}

	return builds
}

func insertServiceInfo(db *sql.DB, info serviceInfo) {
	statement := `
	INSERT INTO services(name, path, args)
	VALUES(?, ?, ?)`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	_, err := db.Exec(statement, info.Name, info.AbsolutePath, info.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteServiceInfo(db *sql.DB, name string) {
	statement := `
	DELETE FROM services
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

func getServiceInfo(db *sql.DB, name string) *serviceInfo {
	statement := `
	SELECT * FROM services s
	WHERE s.name == ?`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	service := serviceInfo{}
	results := db.QueryRow(statement, name)

	if err := results.Scan(&service.ID,
		&service.Name,
		&service.AbsolutePath,
		&service.Args); err != nil {

		fmt.Println(err)
	}

	return &service
}

func getServices(db *sql.DB) []*serviceInfo {
	statement := `
	SELECT * FROM services`

	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err)
	}

	services := make([]*serviceInfo, 0)

	for rows.Next() {
		service := serviceInfo{}
		if err := rows.Scan(&service.ID,
			&service.Name,
			&service.AbsolutePath,
			&service.Args); err != nil {

			fmt.Println(err)
		}
		services = append(services, &service)
	}

	return services
}
