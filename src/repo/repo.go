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
	return Repository{
		Connection: db,
	}
}

// GetSongPath retrieves the song path associated with a given word
func (r *Repository) GetSongPath(word string) string {
	var songPath string
	err := r.Connection.QueryRow("SELECT song_path FROM songs WHERE word = $1", word).Scan(&songPath)
	if err != nil {
		return ""
	}
	return songPath
}

// AddSong adds a new song to the database
func (r *Repository) AddSong(word, songPath string) error {
	_, err := r.Connection.Exec("INSERT into songs (word, song_path) VALUES ($1,$2)", word, songPath)
	return err
}
