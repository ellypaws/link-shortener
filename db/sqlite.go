package db

// use sqlite https://modernc.org/sqlite/

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	dbFile              string = "db.sqlite"
	getCurrentMigration string = `PRAGMA user_version;`
	setCurrentMigration string = `PRAGMA user_version = ?;`
	setForeignKeyCheck  string = `PRAGMA foreign_keys = ON;`
	setBusyTimeOut      string = `PRAGMA busy_timeout = 128;`
)

type Sqlite struct{ *sql.DB }

type migration struct {
	migrationName  string
	migrationQuery string
}

var migrations = []migration{
	{migrationName: "create users table", migrationQuery: createUserTable},
	{migrationName: "create links table", migrationQuery: createLinkTable},
}

// sql statements
const (
	createUserTable = `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY UNIQUE NOT NULL,
		email TEXT,
		password TEXT,
		role INTEGER DEFAULT 1
	)
	`

	createLinkTable = `
	CREATE TABLE IF NOT EXISTS links (
	    short TEXT PRIMARY KEY UNIQUE NOT NULL,
	    original TEXT,
	    owner TEXT,
	    date_created TEXT,
	    date_expired TEXT
	)
	`
)

// New creates a new Sqlite database connection
// Use context to pass in the filename of the database
//
//	context.WithValue(context.Background(), "filename", "audits.sqlite")
//
// Alternatively, use context to pass in ":memory:" to create an in-memory database
func New(ctx context.Context) (*Sqlite, error) {
	var fileName string
	if ctx.Value(":memory:") != nil {
		fileName = ":memory:"
	}

	if s, ok := ctx.Value("filename").(string); ok {
		fileName = s
	}

	if fileName == "" {
		var err error
		fileName, err = DBFilename()
		if err != nil {
			return nil, err
		}

		err = touchDBFile(fileName)
		if err != nil {
			return nil, errors.New("failed to create db file")
		}
	}

	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(setForeignKeyCheck)
	if err != nil {
		return nil, errors.New("failed to enable foreign key checks")
	}

	_, err = db.Exec(setBusyTimeOut)
	if err != nil {
		return nil, errors.New("failed to set busy timeout")
	}

	err = migrate(ctx, db)
	if err != nil {
		return nil, err
	}

	return &Sqlite{DB: db}, nil
}

func DBFilename() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dir + "/" + dbFile, nil
}

func touchDBFile(filename string) error {
	_, err := os.Stat(filename)
	if errors.Is(err, fs.ErrNotExist) {
		file, createErr := os.Create(filename)
		if createErr != nil {
			return createErr
		}

		closeErr := file.Close()
		if closeErr != nil {
			return closeErr
		}
	}

	return nil
}

func migrate(ctx context.Context, db *sql.DB) error {
	var currentMigration int

	row := db.QueryRowContext(ctx, getCurrentMigration)

	err := row.Scan(&currentMigration)
	if err != nil {
		return err
	}

	requiredMigration := len(migrations)

	log.Printf("Current DB version: %v, required DB version: %v\n", currentMigration, requiredMigration)

	if currentMigration < requiredMigration {
		for migrationNum := currentMigration + 1; migrationNum <= requiredMigration; migrationNum++ {
			err = execMigration(ctx, db, migrationNum)
			if err != nil {
				log.Printf("Error running migration %v '%v'\n", migrationNum, migrations[migrationNum-1].migrationName)

				return err
			}
		}
	}

	return nil
}

func execMigration(ctx context.Context, db *sql.DB, migrationNum int) error {
	log.Printf("Running migration %v '%v'\n", migrationNum, migrations[migrationNum-1].migrationName)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, migrations[migrationNum-1].migrationQuery)
	if err != nil {
		return err
	}

	setQuery := strings.Replace(setCurrentMigration, "?", strconv.Itoa(migrationNum), 1)

	_, err = tx.ExecContext(ctx, setQuery)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

var nilDatabase = errors.New("database error")

const timeout = 15 * time.Second

func Error(db *Sqlite) error {
	if db == nil {
		return nilDatabase
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return db.PingContext(ctx)
}

func (db *Sqlite) Wait() {
	var timeout int
	res := db.QueryRow(`PRAGMA busy_timeout;`)
	err := res.Scan(&timeout)
	if err != nil {
		log.Printf("failed to get busy_timeout: %v\n", err)
	}

	sleep := 500 * time.Millisecond
	if timeout > 0 {
		sleep = time.Duration(timeout) * time.Millisecond
	}

	for {
		if err := db.Ping(); err == nil {
			return
		}
		time.Sleep(sleep)
	}
}
