package main

import (
	"fmt"
	"log"
	"os"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	repoURL := fmt.Sprintf("https://%s@github.com/RickDeb2004/VR-Security-Assignment", token)
	localPath := "./temp-repository"

	var repo *git.Repository
	var err error

	// Check if the directory already exists
	if _, err := os.Stat(localPath); !os.IsNotExist(err) {
		fmt.Println("Repository already cloned. Pulling latest changes...")
		repo, err = git.PlainOpen(localPath)
		if err != nil {
			log.Fatalf("Error opening repo: %v", err)
		}

		// Pull the latest changes
		worktree, err := repo.Worktree()
		if err != nil {
			log.Fatalf("Error getting worktree: %v", err)
		}

		err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			log.Fatalf("Error pulling changes: %v", err)
		}

		fmt.Println("Repository updated successfully!")
	} else {
		// Clone the repository if it doesn't exist
		fmt.Println("Cloning the repository...")
		repo, err = git.PlainClone(localPath, false, &git.CloneOptions{
			URL:           repoURL,
			ReferenceName: plumbing.NewBranchReferenceName("main"),
			Progress:      os.Stdout,
		})
		if err != nil {
			log.Fatalf("Error cloning repo: %v", err)
		}

		fmt.Println("Repository cloned successfully!")
	}

	// Make a simple change to the README file
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

	fmt.Println("Changes made successfully!")

	// Commit the changes
	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Error getting worktree: %v", err)
	}

	_, err = worktree.Add("README.md")
	if err != nil {
		log.Fatalf("Error staging file: %v", err)
	}

	commit, err := worktree.Commit("Updated README.md with a new contribution", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Debanjan Mukherjee",
			Email: "debanjanrick04@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Fatalf("Error committing changes: %v", err)
	}

	fmt.Printf("Commit created: %s\n", commit)

	// Push the changes to the remote repository
	auth := &http.BasicAuth{
		Username: "RickDeb2004", // The username can be any non-empty string
		Password: token,
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		log.Fatalf("Error pushing changes: %v", err)
	}

	fmt.Println("Changes pushed successfully!")
}
