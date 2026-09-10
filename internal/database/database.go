package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	query := `
		CREATE TABLE IF NOT EXISTS thumbnails (
			video_id TEXT PRIMARY KEY,
			data BLOB NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	if _, err = db.Exec(query); err != nil {
		db.Close()
		return nil, fmt.Errorf("create thumbnails table: %w", err)
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetThumbnail(videoID string) ([]byte, error) {
	var data []byte

	err := d.db.QueryRow(
		"SELECT data FROM thumbnails WHERE video_id = ?",
		videoID,
	).Scan(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (d *Database) SaveThumbnail(videoID string, data []byte) error {
	_, err := d.db.Exec(
		"INSERT INTO thumbnails (video_id, data) VALUES (?, ?)",
		videoID,
		data,
	)

	if err != nil {
		return err
	}

	return nil
}
