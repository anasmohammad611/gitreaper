package cli

import (
	"fmt"
	"github.com/anasmohammad611/gitreaper/internal/git"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean merged branches from your repository",
	Long: `Clean removes branches that have been merged into main branches.
	
By default, it will:
- Find branches merged into main/staging/dev/development/production/preproduction
- Show you what would be deleted
- Ask for confirmation before deletion`,
	RunE: runClean,
}

func runClean(cmd *cobra.Command, args []string) error {
	fmt.Println("🧹 GitReaper Clean - Starting repository cleanup...")

	// Step 1: Open Git repository
	repo, err := git.NewRepository()
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	fmt.Printf("📁 Repository: %s\n", repo.GetPath())

	// Step 2: Find main branches
	mainBranches, err := repo.GetMainBranches()
	if err != nil {
		return fmt.Errorf("failed to find main branches: %w", err)
	}

	fmt.Printf("🎯 Found main branches: %v\n", mainBranches)

	// Step 3: Find merged branches
	fmt.Println("\n🔍 Searching for merged branches...")
	mergedBranches, err := repo.GetMergedBranches(mainBranches)
	if err != nil {
		return fmt.Errorf("failed to find merged branches: %w", err)
	}

	// Step 4: Display results
	if len(mergedBranches) == 0 {
		fmt.Println("✅ No merged branches found - your repository is already clean!")
		return nil
	}

	fmt.Printf("\n📋 Found %d merged branch(es):\n\n", len(mergedBranches))
	for _, branch := range mergedBranches {
		fmt.Printf("  • %s (merged into %s) - %s by %s\n",
			branch.Name,
			branch.MergedInto,
			branch.LastCommit,
			branch.LastAuthor)
	}

	// TODO: Ask for confirmation
	// TODO: Delete branches

	return nil
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
