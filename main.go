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
		repo.Commit(os.Args[2:])
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
