package main

import (
	"fmt"
	"log"
	"os"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func main() {
	repoURL := "https://github.com/RickDeb2004/Auto_demo_contri" // Replace with your repo
	localPath := "./temp-repo"

	// Clone the repository
	fmt.Println("Cloning the repository...")
	_, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		log.Fatalf("Error cloning repo: %v", err)
	}

	// Open the repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		log.Fatalf("Error opening repo: %v", err)
	}

	// Create a new file in the repo
	filePath := localPath + "/README.md"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\nContribution made on %s", time.Now().Format(time.RFC1123)))
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	// Commit the changes
	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Error getting worktree: %v", err)
	}

	_, err = worktree.Add("README.md")
	if err != nil {
		log.Fatalf("Error adding file to worktree: %v", err)
	}

	_, err = worktree.Commit("Automated commit by Go bot", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "GitHub Bot",
			Email: "bot@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Fatalf("Error committing changes: %v", err)
	}

	fmt.Println("Changes committed!")

	// Push the changes
	fmt.Println("Pushing to GitHub...")
	err = repo.Push(&git.PushOptions{})
	if err != nil {
		log.Fatalf("Error pushing to GitHub: %v", err)
	}

	fmt.Println("Pushed successfully!")
}
