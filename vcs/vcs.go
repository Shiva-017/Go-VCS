package vcs

import (
	"fmt"
	"go-vcs/merkle"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Repository stores commit history
type Repository struct{}

// NewRepository initializes a new VCS repository
func NewRepository() *Repository {
	InitDB()
	return &Repository{}
}

// Add reads and returns a map of file path to content
func (repo *Repository) Add(files []string) map[string]string {
	fileData := make(map[string]string)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("error reading file %s: %v\n", file, err)
			continue
		}

		fileData[file] = string(content)
		fmt.Printf("staged: %s\n", file)
	}

	return fileData
}

// Commit creates a Merkle-tree snapshot of the repository
func (repo *Repository) Commit(files []string, message string) {
	fileData := repo.Add(files)
	var hashes []string

	for _, content := range fileData {
		hashes = append(hashes, merkle.ComputeHash(content))
	}

	root := merkle.BuildMerkleTree(hashes)
	if root == nil {
		fmt.Println("nothing to commit")
		return
	}

	commitID := merkle.ComputeHash(time.Now().String())
	timestamp := time.Now().Format(time.RFC3339)

	SaveCommit(commitID, timestamp, root.Hash, fileData, message)
	fmt.Println("Commit successful! Root hash:", root.Hash)
}

// History prints all commits in chronological order
func (repo *Repository) History() {
	commitIDs, messages := GetCommitHistory()

	fmt.Println("📜 Commit History:")
	if len(commitIDs) == 0 {
		fmt.Println("⚠️ No commits found.")
		return
	}

	for i := 0; i < len(commitIDs); i++ {
		fmt.Printf("  %s  %s\n", commitIDs[i], messages[i])
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

func GetAllFiles(directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Only include files that are not directories and not hidden
		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
