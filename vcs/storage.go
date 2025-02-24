package vcs

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "vcs.db"

// initializes the database
func InitDB() *sql.DB {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			fmt.Println("Error creating database:", err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Println("Error opening database:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS commits (
		id TEXT PRIMARY KEY,
		timestamp TEXT,
		root_hash TEXT
	);
	CREATE TABLE IF NOT EXISTS files (
		commit_id TEXT,
		filename TEXT,
		content TEXT,
		FOREIGN KEY(commit_id) REFERENCES commits(id)
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Error creating tables:", err)
	}

	return db
}

// SaveCommit stores a commit in the database
func SaveCommit(commitID, timestamp, rootHash string, files map[string]string) {
	db := InitDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO commits (id, timestamp, root_hash) VALUES (?, ?, ?)", commitID, timestamp, rootHash)
	if err != nil {
		fmt.Println("Error saving commit:", err)
	}

	for filename, content := range files {
		_, err = db.Exec("INSERT INTO files (commit_id, filename, content) VALUES (?, ?, ?)", commitID, filename, content)
		if err != nil {
			fmt.Println("Error saving file:", err)
		}
	}
}

// GetCommit retrieves a commit and its associated files
func GetCommit(commitID string) (string, map[string]string) {
	db := InitDB()
	defer db.Close()

	var rootHash string
	err := db.QueryRow("SELECT root_hash FROM commits WHERE id = ?", commitID).Scan(&rootHash)
	if err != nil {
		fmt.Println("❌ Error: Commit not found. Make sure you are using the full commit ID.")
		return "", nil
	}

	rows, err := db.Query("SELECT filename, content FROM files WHERE commit_id = ?", commitID)
	if err != nil {
		fmt.Println("❌ Error retrieving files for commit:", commitID)
		return "", nil
	}
	defer rows.Close()

	files := make(map[string]string)
	for rows.Next() {
		var filename, content string
		if err := rows.Scan(&filename, &content); err != nil {
			fmt.Println("⚠️ Error scanning file data:", err)
			continue
		}
		files[filename] = content
	}

	if len(files) == 0 {
		fmt.Println("⚠️ No files found for this commit.")
	}

	return rootHash, files
}

// GetCommitHistory returns all commit IDs
func GetCommitHistory() []string {
	db := InitDB()
	defer db.Close()

	rows, err := db.Query("SELECT id FROM commits")
	if err != nil {
		fmt.Println("Error retrieving history:", err)
		return nil
	}
	defer rows.Close()

	var commitIDs []string
	for rows.Next() {
		var commitID string
		rows.Scan(&commitID)
		commitIDs = append(commitIDs, commitID)
	}

	return commitIDs
}
