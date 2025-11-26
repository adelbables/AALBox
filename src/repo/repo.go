// Package repo provides functionalities to interact with a database
// for managing songs associated with tags.package repo
package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Repository struct to manage database interactions
type Repository struct {
	Connection *sql.DB
}

// NewRepository initializes a new Repository with a database connection
func NewRepository() Repository {
	connStr := "user=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Failed to open database connection:", err)
		panic(err)
	}

	// It's a good practice to test the database connection
	if err := db.Ping(); err != nil {
		fmt.Println("Failed to establish database connection:", err)
		panic(err)
	}

	tableCreationQuery := `CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    tag_id TEXT NOT NULL,
    song_path TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_songs_tag_id ON songs(tag_id);`

	if _, err := db.Exec(tableCreationQuery); err != nil {
		fmt.Println("Failed to create table:", err)
	}

	return Repository{
		Connection: db,
	}
}

// GetSongPath retrieves the song path associated with a given word
func (r *Repository) GetSongPath(tagID string) string {
	var songPath string
	err := r.Connection.QueryRow(
		"SELECT song_path FROM songs WHERE tag_id = $1 ORDER BY random() LIMIT 1",
		tagID,
	).Scan(&songPath)
	if err != nil {
		return ""
	}
	return songPath
}

// AddSong adds a new song to the database
func (r *Repository) AddSong(tagID, songPath string) error {
	_, err := r.Connection.Exec(
		"INSERT INTO songs (tag_id, song_path) VALUES ($1, $2)",
		tagID, songPath,
	)
	return err
}

func (r *Repository) RemoveSong(tagID, songPath string) error {
	_, err := r.Connection.Exec(
		"DELETE FROM songs WHERE tag_id = $1 AND song_path = $2",
		tagID, songPath,
	)
	return err
}

func (r *Repository) RemoveAllSongs(tagID string) error {
	_, err := r.Connection.Exec("DELETE FROM songs WHERE tag_id =$1", tagID)
	return err
}
