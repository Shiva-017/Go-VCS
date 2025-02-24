package vcs

import (
	"fmt"
	"go-vcs/merkle"
	"os"
	"time"
)

// Repository stores commit history
type Repository struct{}

// NewRepository initializes a repository
func NewRepository() *Repository {
	InitDB()
	return &Repository{}
}

// Add reads file contents
func (repo *Repository) Add(files []string) map[string]string {
	fileData := make(map[string]string)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", file)
			continue
		}

		fileData[file] = string(content)
		fmt.Println("Added:", file)
	}

	return fileData
}

// Commit creates a snapshot
func (repo *Repository) Commit(files []string) {
	fileData := repo.Add(files)
	var hashes []string

	for _, content := range fileData {
		hashes = append(hashes, merkle.ComputeHash(content))
	}

	root := merkle.BuildMerkleTree(hashes)
	if root == nil {
		fmt.Println("No files to commit.")
		return
	}

	commitID := merkle.ComputeHash(time.Now().String())
	timestamp := time.Now().Format(time.RFC3339)

	SaveCommit(commitID, timestamp, root.Hash, fileData)
	fmt.Println("Commit successful! Root hash:", root.Hash)
}

// History shows commit history
func (repo *Repository) History() {
	commitIDs := GetCommitHistory()
	fmt.Println("📜 Commit History:")
	for _, id := range commitIDs {
		fmt.Println("🔹 Commit ID:", id) // Print full commit ID
	}
}

// Revert restores files
func (repo *Repository) Revert(commitID string) {
	rootHash, files := GetCommit(commitID)

	if rootHash == "" || files == nil {
		fmt.Println("❌ Error: Commit not found. Ensure you are using the full commit ID from history.")
		return
	}

	fmt.Println("🔄 Reverting to commit:", commitID)
	for filename, content := range files {
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			fmt.Println("⚠️ Error restoring", filename, ":", err)
			continue
		}
		fmt.Println("✅ Restored:", filename)
	}

	fmt.Println("✔️ Revert successful!")
}
