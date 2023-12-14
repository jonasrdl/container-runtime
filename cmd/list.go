package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all running containers",
	RunE:  list,
}

var rootPath = "/var/lib/container-runtime"

func list(_ *cobra.Command, _ []string) error {
	containerDir, err := os.ReadDir(rootPath)
	if err != nil {
		return fmt.Errorf("error reading root path: %v", err)
	}

	var validContainers []string

	for _, container := range containerDir {
		if container.IsDir() && isSHA256Hash(container.Name()) {
			validContainers = append(validContainers, container.Name())
		}
	}

	if len(validContainers) == 0 {
		fmt.Println("No containers are running")
		return nil
	}

	// print out containers
	for _, container := range validContainers {
		fmt.Println(container)
	}
	return nil
}

// isSHA256Hash checks if a string is a valid SHA256 hash
func isSHA256Hash(hash string) bool {
	return regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(hash)
}

func init() {
	rootCmd.AddCommand(listCmd)
}
