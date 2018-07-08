package main

import (
	"database/sql"
	"os"
	"testing"
)

func TestBuildInfoInsertion(t *testing.T) {
	const test1 = "test_1"
	const test2 = "test_2"
	const test3 = "test_3"

	db, err := sql.Open("sqlite3", "test.db")
	defer func() {
		db.Close()
		os.Remove("test.db")
	}()

	if err != nil {
		t.Error(err)
	}

	statement := `
	CREATE TABLE IF NOT EXISTS builds(
		id INTEGER NOT NULL PRIMARY KEY,
		name STRING NOT NULL UNIQUE,
		path STRING NOT NULL,
		branch STRING);`

	_, err = db.Exec(statement)
	if err != nil {
		t.Error(err)
	}

	in := buildInfo{
		Name:         test1,
		AbsolutePath: test2,
		Branch:       test3,
	}
	insertBuildInfo(db, in)

	out := getBuildInfo(db, test1)
	if in.Name != out.Name ||
		in.AbsolutePath != out.AbsolutePath ||
		in.Branch != out.Branch {

		t.Error("Input and output did not match!")
	}

}
