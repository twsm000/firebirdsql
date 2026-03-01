package test

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/nakagami/firebirdsql"
)

var (
	//go:embed sql/empty_database.fdb
	emptyFirebirdDatabase []byte

	//go:embed sql/*.sql
	scripts embed.FS
)

func EmptyFirebirdDatabase() io.Reader {
	return bytes.NewReader(emptyFirebirdDatabase)
}

type DBConfig struct {
	Driver      string
	Host        string
	User        string
	Database    string
	Password    string
	Charset     string
	MaxOpenConn int
	MaxIdleConn int
	MaxIdleTime time.Duration
}

func executeScripts(db *sql.DB) error {
	var err error
	err = errors.Join(err, executeScript(db, "TEST_TABLE.sql"))
	return err
}

func executeScript(db *sql.DB, scriptName string) error {
	data, err := scripts.ReadFile("sql/" + scriptName)
	if err != nil {
		return fmt.Errorf("failed to read script: %s: %w", scriptName, err)
	}
	queries := strings.Split(string(data), "--;")
	for _, sql := range queries {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		sql = strings.TrimSuffix(sql, ";")
		if _, err := db.Exec(sql); err != nil {
			return fmt.Errorf("failed to execute script: %s: %w - query = %s", scriptName, err, sql)
		}
	}
	return nil
}

func CreateTempDB(dbname string) (db *sql.DB, cleanup func(), err error) {
	dst, err := os.CreateTemp(os.TempDir(), dbname+"*.fdb")
	if err != nil {
		err = fmt.Errorf("failed to create db: %v - %w", dst, err)
		return
	}
	tmpDB := dst.Name()
	cleanup = func() {
		_ = dst.Close()
		if db != nil {
			log.Println("closing database:", tmpDB)
			if err := db.Close(); err != nil {
				log.Println("failed to close database:", tmpDB, err)
			}
		}

		log.Println("removing db file:", tmpDB)
		if err := os.Remove(tmpDB); err != nil {
			log.Println("failed to remove file:", tmpDB, err)
		}
	}

	src := EmptyFirebirdDatabase()
	if _, err = io.Copy(dst, src); err != nil {
		err = fmt.Errorf("failed to copy db file: %v - %w", dst, err)
		return
	}
	_ = dst.Close()

	dbConfig := DBConfig{
		Driver:      "firebirdsql",
		Host:        os.Getenv("HOSTNAME"),
		User:        "SYSDBA",
		Database:    tmpDB,
		Password:    "masterkey",
		Charset:     "WIN1252",
		MaxOpenConn: 1,
		MaxIdleConn: 1,
		MaxIdleTime: 10 * time.Minute,
	}

	db, err = NewConnection(dbConfig)
	if err != nil {
		err = fmt.Errorf("failed to connect to db: %v - %w", tmpDB, err)
		return
	}
	log.Println("db created:", tmpDB)

	if err = executeScripts(db); err != nil {
		err = fmt.Errorf("failed to execute scripts: %w", err)
		return
	}
	return
}

func NewConnection(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=%s",
		config.User, config.Password, config.Host, config.Database, config.Charset,
	)
	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConn)
	db.SetMaxIdleConns(config.MaxIdleConn)
	db.SetConnMaxIdleTime(config.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
