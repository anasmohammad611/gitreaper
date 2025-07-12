package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"os"
)

// Repository wraps go-git functionality
type Repository struct {
	repo *git.Repository
	path string
}

// NewRepository opens a Git repository in the current directory
func NewRepository() (*Repository, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	// Find the Git repository (searches up the directory tree)
	repo, err := git.PlainOpenWithOptions(cwd, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return nil, fmt.Errorf("not a Git repository (or any parent directory)")
		}
		return nil, fmt.Errorf("failed to open Git repository: %w", err)
	}

	// Get the repository root path
	workTree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	return &Repository{
		repo: repo,
		path: workTree.Filesystem.Root(),
	}, nil
}

// GetPath returns the repository root path
func (r *Repository) GetPath() string {
	return r.path
}

// GetMainBranches finds common main branches that exist in the repository
func (r *Repository) GetMainBranches() ([]string, error) {
	commonMainBranches := []string{
		"main",
		"master",
		"staging",
		"dev",
		"development",
		"production",
		"preproduction",
	}

	var existingMainBranches []string

	// Get all local branches
	branches, err := r.repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	// Check which main branches exist
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		branchName := ref.Name().Short()

		// Check if this branch is one of our main branches
		for _, mainBranch := range commonMainBranches {
			if branchName == mainBranch {
				existingMainBranches = append(existingMainBranches, branchName)
				log.Debug().Str("branch", branchName).Msg("Found main branch")
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	if len(existingMainBranches) == 0 {
		return nil, fmt.Errorf("no main branches found (looking for: %v)", commonMainBranches)
	}

	return existingMainBranches, nil
}

// BranchInfo contains information about a branch
type BranchInfo struct {
	Name       string
	MergedInto string // Which main branch it's merged into
	LastCommit string // Last commit hash (short)
	LastAuthor string // Last commit author
}

// GetMergedBranches finds branches that are merged into main branches
func (r *Repository) GetMergedBranches(mainBranches []string) ([]BranchInfo, error) {
	var mergedBranches []BranchInfo

	// Get all local branches
	branches, err := r.repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		branchName := ref.Name().Short()

		// Skip main branches themselves
		isMainBranch := false
		for _, mainBranch := range mainBranches {
			if branchName == mainBranch {
				isMainBranch = true
				break
			}
		}
		if isMainBranch {
			return nil
		}

		// Check if this branch is merged into any main branch
		for _, mainBranch := range mainBranches {
			isMerged, err := r.isBranchMerged(branchName, mainBranch)
			if err != nil {
				log.Warn().Err(err).
					Str("branch", branchName).
					Str("main", mainBranch).
					Msg("Failed to check if branch is merged")
				continue
			}

			if isMerged {
				// Get branch info
				branchInfo, err := r.getBranchInfo(branchName, mainBranch)
				if err != nil {
					log.Warn().Err(err).Str("branch", branchName).Msg("Failed to get branch info")
					continue
				}

				mergedBranches = append(mergedBranches, branchInfo)
				break // Don't check other main branches for this branch
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	return mergedBranches, nil
}

// isBranchMerged checks if a branch is merged into a main branch
func (r *Repository) isBranchMerged(branchName, mainBranch string) (bool, error) {
	// Get the commit of the branch
	branchRef, err := r.repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err != nil {
		return false, fmt.Errorf("failed to get branch reference: %w", err)
	}

	// Get the commit of the main branch
	mainRef, err := r.repo.Reference(plumbing.NewBranchReferenceName(mainBranch), true)
	if err != nil {
		return false, fmt.Errorf("failed to get main branch reference: %w", err)
	}

	// Check if branch commit is reachable from main branch
	// This means the branch has been merged
	branchCommit, err := r.repo.CommitObject(branchRef.Hash())
	if err != nil {
		return false, fmt.Errorf("failed to get branch commit: %w", err)
	}

	mainCommit, err := r.repo.CommitObject(mainRef.Hash())
	if err != nil {
		return false, fmt.Errorf("failed to get main commit: %w", err)
	}

	// Use git's merge-base logic
	isAncestor, err := branchCommit.IsAncestor(mainCommit)
	if err != nil {
		return false, fmt.Errorf("failed to check ancestry: %w", err)
	}

	return isAncestor, nil
}

// getBranchInfo gets detailed information about a branch
func (r *Repository) getBranchInfo(branchName, mergedInto string) (BranchInfo, error) {
	// Get branch reference
	ref, err := r.repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err != nil {
		return BranchInfo{}, fmt.Errorf("failed to get branch reference: %w", err)
	}

	// Get last commit
	commit, err := r.repo.CommitObject(ref.Hash())
	if err != nil {
		return BranchInfo{}, fmt.Errorf("failed to get commit: %w", err)
	}

	return BranchInfo{
		Name:       branchName,
		MergedInto: mergedInto,
		LastCommit: ref.Hash().String()[:8], // Short hash
		LastAuthor: commit.Author.Name,
	}, nil
}

// DeleteBranch deletes a local branch
func (r *Repository) DeleteBranch(branchName string) error {
	// Get branch reference
	refName := plumbing.NewBranchReferenceName(branchName)

	// Delete the reference
	err := r.repo.Storer.RemoveReference(refName)
	if err != nil {
		return fmt.Errorf("failed to delete branch %s: %w", branchName, err)
	}

	return nil
}

// DeleteBranches deletes multiple branches and returns results
func (r *Repository) DeleteBranches(branches []BranchInfo) ([]string, []error) {
	var deleted []string
	var errors []error

	for _, branch := range branches {
		err := r.DeleteBranch(branch.Name)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to delete %s: %w", branch.Name, err))
		} else {
			deleted = append(deleted, branch.Name)
		}
	}

	return deleted, errors
}
