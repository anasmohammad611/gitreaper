package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConfirmDeletion asks the user to confirm branch deletion
func ConfirmDeletion(branchCount int) (bool, error) {
	fmt.Printf("\n❓ Do you want to delete %d merged branch(es)? [y/N]: ", branchCount)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}

	// Clean the response
	response = strings.TrimSpace(strings.ToLower(response))

	// Accept 'y', 'yes' as confirmation
	return response == "y" || response == "yes", nil
}
