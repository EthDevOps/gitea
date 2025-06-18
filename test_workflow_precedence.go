package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code.gitea.io/gitea/modules/actions"
	"code.gitea.io/gitea/modules/git"
)

func main() {
	// Test the workflow precedence behavior
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run test_workflow_precedence.go <repo_path>")
		os.Exit(1)
	}

	repoPath := os.Args[1]
	
	// Open repository
	repo, err := git.OpenRepository(git.DefaultContext, repoPath)
	if err != nil {
		fmt.Printf("Error opening repository: %v\n", err)
		os.Exit(1)
	}
	defer repo.Close()

	// Get the HEAD commit
	commit, err := repo.GetCommit("HEAD")
	if err != nil {
		fmt.Printf("Error getting HEAD commit: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Testing workflow precedence for commit: %s\n", commit.ID.String())
	fmt.Printf("Commit message: %s\n", strings.Split(commit.CommitMessage, "\n")[0])
	fmt.Println()

	// Test ListWorkflows function
	folder, entries, err := actions.ListWorkflows(commit)
	if err != nil {
		fmt.Printf("Error listing workflows: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ListWorkflows result:\n")
	fmt.Printf("  Folder: %s\n", folder)
	fmt.Printf("  Number of entries: %d\n", len(entries))
	
	if len(entries) > 0 {
		fmt.Printf("  Workflow files found:\n")
		for i, entry := range entries {
			fmt.Printf("    %d. %s\n", i+1, entry.Name())
		}
	}
	fmt.Println()

	// Let's also manually check what exists in both directories
	fmt.Printf("Manual directory check:\n")
	
	// Check .gitea/workflows
	giteaTree, err := commit.SubTree(".gitea/workflows")
	if err != nil {
		fmt.Printf("  .gitea/workflows: Not found (%v)\n", err)
	} else {
		giteaEntries, err := giteaTree.ListEntriesRecursiveFast()
		if err != nil {
			fmt.Printf("  .gitea/workflows: Error listing (%v)\n", err)
		} else {
			fmt.Printf("  .gitea/workflows: Found %d files\n", len(giteaEntries))
			for _, entry := range giteaEntries {
				if strings.HasSuffix(entry.Name(), ".yml") || strings.HasSuffix(entry.Name(), ".yaml") {
					fmt.Printf("    - %s\n", entry.Name())
				}
			}
		}
	}

	// Check .github/workflows
	githubTree, err := commit.SubTree(".github/workflows")
	if err != nil {
		fmt.Printf("  .github/workflows: Not found (%v)\n", err)
	} else {
		githubEntries, err := githubTree.ListEntriesRecursiveFast()
		if err != nil {
			fmt.Printf("  .github/workflows: Error listing (%v)\n", err)
		} else {
			fmt.Printf("  .github/workflows: Found %d files\n", len(githubEntries))
			for _, entry := range githubEntries {
				if strings.HasSuffix(entry.Name(), ".yml") || strings.HasSuffix(entry.Name(), ".yaml") {
					fmt.Printf("    - %s\n", entry.Name())
				}
			}
		}
	}
}