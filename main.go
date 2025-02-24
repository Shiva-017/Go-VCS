package main

import (
	"fmt"
	"go-vcs/vcs"
	"os"
)

func main() {
	repo := vcs.NewRepository()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go-vcs <command> [args]")
		return
	}

	command := os.Args[1]

	switch command {
	case "commit":
		files, err := vcs.GetAllFiles("./Repository")
		if err != nil {
			fmt.Println("❌ Error getting files:", err)
			return
		}

		if len(files) == 0 {
			fmt.Println("⚠️ No files to commit.")
			return
		}

		var commitMessage string
		if len(os.Args) > 2 {
			commitMessage = os.Args[2]
		} else {
			commitMessage = "Default commit message"
		}

		repo.Commit(files, commitMessage)
	case "history":
		repo.History()
	case "revert":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go-vcs revert <commit-id>")
			return
		}
		commitID := os.Args[2]
		repo.Revert(commitID)
	default:
		fmt.Println("Unknown command:", command)
	}
}
