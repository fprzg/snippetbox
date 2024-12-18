package models

import (
	"database/sql"
	"os"
	"testing"
)

func dbExecuteScript(db *sql.DB, filepath string) error {
	script, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(script))
	if err != nil {
		return err
	}

	return nil
}

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	err = dbExecuteScript(db, "./db-scripts/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = dbExecuteScript(db, "./db-scripts/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return db
}
