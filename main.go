// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"time"

// 	git "gopkg.in/src-d/go-git.v4"
// 	"gopkg.in/src-d/go-git.v4/plumbing/object"
// )

// func main() {
// 	// Get the current working directory instead of creating a new one
// 	workDir, err := os.Getwd()
// 	if err != nil {
// 		log.Fatalf("Error getting working directory: %v", err)
// 	}

// 	// Open the existing repository instead of cloning
// 	repo, err := git.PlainOpen(workDir)
// 	if err != nil {
// 		log.Fatalf("Error opening repository: %v", err)
// 	}

// 	// Make changes to README.md
// 	filePath := workDir + "/README.md"
// 	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
// 	if err != nil {
// 		log.Fatalf("Error opening file: %v", err)
// 	}
// 	defer file.Close()

// 	_, err = file.WriteString(fmt.Sprintf("\nContribution made on %s", time.Now().Format(time.RFC1123)))
// 	if err != nil {
// 		log.Fatalf("Error writing to file: %v", err)
// 	}
// 	fmt.Println("Changes made successfully!")

// 	// Commit the changes
// 	worktree, err := repo.Worktree()
// 	if err != nil {
// 		log.Fatalf("Error getting worktree: %v", err)
// 	}

// 	_, err = worktree.Add("README.md")
// 	if err != nil {
// 		log.Fatalf("Error staging file: %v", err)
// 	}

// 	commit, err := worktree.Commit("Updated README.md with a new contribution", &git.CommitOptions{
// 		Author: &object.Signature{
// 			Name:  "github-actions[bot]",
// 			Email: "github-actions[bot]@users.noreply.github.com",
// 			When:  time.Now(),
// 		},
// 	})
// 	if err != nil {
// 		log.Fatalf("Error committing changes: %v", err)
// 	}
// 	fmt.Printf("Commit created: %s\n", commit)

//		// Push using git command line instead of go-git
//		pushCmd := "git push origin HEAD:main"
//		output, err := exec.Command("sh", "-c", pushCmd).CombinedOutput()
//		if err != nil {
//			log.Fatalf("Error pushing changes: %v\nOutput: %s", err, string(output))
//		}
//		fmt.Println("Changes pushed successfully!")
//	}
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Repository struct to hold repo information
type Repository struct {
	Name string
	URL  string
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	// List of target repositories
	repositories := []Repository{
		{
			Name: "VR-Security-Assignment",
			URL:  "https://github.com/RickDeb2004/VR-Security-Assignment",
		},
		{
			Name: "VR",
			URL:  "https://github.com/RickDeb2004/vr",
		},
		 

		// Add more repositories here as needed
	}

	baseDir := "./repos"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Fatalf("Error creating base directory: %v", err)
	}

	for _, repo := range repositories {
		processRepository(repo, baseDir, token)
	}
}

func processRepository(repo Repository, baseDir, token string) {
	repoDir := filepath.Join(baseDir, repo.Name)
	repoURL := fmt.Sprintf("https://RickDeb2004:%s@%s", token, repo.URL[8:]) // Convert https://github.com to include token

	// Clone repository
	fmt.Printf("Cloning %s...\n", repo.Name)
	r, err := git.PlainClone(repoDir, false, &git.CloneOptions{
		URL: repoURL,
		Auth: &http.BasicAuth{
			Username: "RickDeb2004",
			Password: token,
		},
	})
	if err != nil {
		log.Printf("Error cloning %s: %v\n", repo.Name, err)
		return
	}

	// Make 20 commits
	for i := 1; i <= 20; i++ {
		if err := makeCommit(r, repoDir, i, repo.Name); err != nil {
			log.Printf("Error making commit %d for %s: %v\n", i, repo.Name, err)
			continue
		}

		// Push after each commit
		if err := pushChanges(repoDir, token); err != nil {
			log.Printf("Error pushing commit %d for %s: %v\n", i, repo.Name, err)
			continue
		}

		// Add delay between commits
		time.Sleep(2 * time.Second)
	}
}

func makeCommit(r *git.Repository, repoDir string, commitNum int, repoName string) error {
	// Update README.md
	readmePath := filepath.Join(repoDir, "README.md")
	file, err := os.OpenFile(readmePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening README: %v", err)
	}
	defer file.Close()

	commitMsg := fmt.Sprintf("\nCommit #%d made on %s for %s\n", 
		commitNum,
		time.Now().Format(time.RFC1123),
		repoName)
	
	if _, err := file.WriteString(commitMsg); err != nil {
		return fmt.Errorf("error writing to README: %v", err)
	}

	// Stage changes
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %v", err)
	}

	_, err = w.Add("README.md")
	if err != nil {
		return fmt.Errorf("error staging file: %v", err)
	}

	// Commit changes
	_, err = w.Commit(fmt.Sprintf("Update #%d: Daily contribution", commitNum), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "RickDeb2004",
			Email: "debanjanrick04@gmail.com",
			When:  time.Now(),
		},
	})
	
	return err
}

func pushChanges(repoDir string, token string) error {
	cmd := exec.Command("git", "push", "origin", "main")
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GIT_ASKPASS=echo"),
		fmt.Sprintf("GIT_USERNAME=RickDeb2004"),
		fmt.Sprintf("GIT_PASSWORD=%s", token))
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("push error: %v, output: %s", err, string(output))
	}
	return nil
}