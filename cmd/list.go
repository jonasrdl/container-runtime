package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all running containers",
	Run:   list,
}

var rootPath = "/var/lib/container-runtime"

func list(_ *cobra.Command, _ []string) {
	containerDir, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println("Error reading root path:", err)
		os.Exit(1)
	}
	var containers []string
	for _, container := range containerDir {
		if container.IsDir() {
			containers = append(containers, container.Name())
		}
	}

	if len(containers) == 0 {
		fmt.Println("No containers are running")
		os.Exit(0)
	}

	// print out containers
	for _, container := range containers {
		fmt.Println(container)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
