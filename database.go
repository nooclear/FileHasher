package FileHasher

import (
	"context"
	"database/sql"
	"errors"
	"log"
	_ "modernc.org/sqlite"
)

var db *sql.DB

type entry struct {
	id        int
	timestamp int64
	filepath  string
	hash      string
}

func initDB(dbPath string) error {
	var err error
	if db, err = sql.Open("sqlite", dbPath); err != nil {
		return err
	}
	if _, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS hashes (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	timestamp INTEGER NOT NULL,
    	filepath TEXT NOT NULL,
    	hash TEXT
		)`,
	); err != nil {
		return err
	}

	return nil
}

func add(e *entry) (int64, error) {
	if result, err := db.ExecContext(
		context.Background(),
		`INSERT INTO hashes (timestamp, filepath, hash) VALUES (?, ?, ?)`, e.timestamp, e.filepath, e.hash,
	); err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func update(e *entry) (int64, error) {
	if result, err := db.ExecContext(
		context.Background(),
		`UPDATE hashes SET timestamp = ?, hash = ? WHERE filepath = ?`, e.timestamp, e.hash, e.filepath,
	); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func del(e *entry) (int64, error) {
	if result, err := db.ExecContext(
		context.Background(),
		`DELETE FROM hashes WHERE filepath = ?`, e.filepath,
	); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func find(e *entry) bool {
	if err := db.QueryRow(
		`SELECT filepath FROM hashes WHERE filepath = ?`, e.filepath,
	).Scan(&e.filepath); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Fatal(err)
		}
		return false
	}
	return true
}

func getAll() ([]*entry, error) {
	var entries []*entry
	rows, err := db.QueryContext(context.Background(), `SELECT timestamp, filepath, hash FROM hashes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var e entry
		if err := rows.Scan(&e.timestamp, &e.filepath, &e.hash); err != nil {
			return nil, err
		}
		entries = append(entries, &e)
	}
	return entries, err
}
