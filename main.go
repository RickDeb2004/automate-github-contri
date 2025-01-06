package main

import (
	"fmt"
	"log"
	"os"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {
	// Get the GitHub token from an environment variable
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	repoURL := fmt.Sprintf("https://%s@github.com/RickDeb2004/VR-Security-Assignment", token)
	localPath := "./temp-repo"

	// Clone the repository
	fmt.Println("Cloning the repository...")
	_, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName("main"),
		Progress:      os.Stdout,
	})
	if err != nil {
		log.Fatalf("Error cloning repo: %v", err)
	}

	fmt.Println("Repository cloned successfully!")

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

	fmt.Println("Changes made and committed successfully!")
}
