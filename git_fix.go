package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting cwd: %v\n", err)
		return
	}

	gitPath := filepath.Join(cwd, ".git")

	// 1. Force removal of .git
	fmt.Printf("Purging %s...\n", gitPath)
	err = os.RemoveAll(gitPath)
	if err != nil {
		fmt.Printf("Force removal failed: %v. Trying alternative...\n", err)
		// Try renaming it if we can't delete it
		err = os.Rename(gitPath, filepath.Join(cwd, ".git_stale"))
		if err != nil {
			fmt.Printf("Rename also failed: %v\n", err)
		}
	}

	// 2. Initialize
	fmt.Println("Initializing new Git repo...")
	out, err := exec.Command("git", "init").CombinedOutput()
	fmt.Printf("%s\n", string(out))
	if err != nil {
		return
	}

	// 3. Add
	fmt.Println("Adding files...")
	out, err = exec.Command("git", "add", ".").CombinedOutput()
	fmt.Printf("%s\n", string(out))
	if err != nil {
		return
	}

	// 4. Commit
	fmt.Println("Committing...")
	out, err = exec.Command("git", "commit", "-m", "Initial commit // Cleaned & Optimized").CombinedOutput()
	fmt.Printf("%s\n", string(out))
	if err != nil {
		return
	}

	// 5. Remote and Push
	fmt.Println("Setting remote and pushing...")
	exec.Command("git", "branch", "-M", "main").Run()
	exec.Command("git", "remote", "add", "origin", "https://github.com/w1hi4/VaultGuard.git").Run()
	out, err = exec.Command("git", "push", "-u", "origin", "main", "--force").CombinedOutput()
	fmt.Printf("%s\n", string(out))
}
