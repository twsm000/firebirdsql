package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	db, cleanup, err := CreateTempDB("test_database")
	if err != nil {
		fmt.Println("failed to create temp db or run scripts:", err)
		os.Exit(1)
	}
	if db == nil {
		fmt.Println("db connection should not be nil")
		os.Exit(1)
	}
	testDB = db

	code := m.Run()
	cleanup()
	os.Exit(code)
}
