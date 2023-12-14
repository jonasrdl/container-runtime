package cmd

import (
	"fmt"
	"os"

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
	var containers []string
	for _, container := range containerDir {
		if container.IsDir() {
			containers = append(containers, container.Name())
		}
	}

	if len(containers) == 0 {
		return fmt.Errorf("no containers are running")
	}

	// print out containers
	for _, container := range containers {
		fmt.Println(container)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
