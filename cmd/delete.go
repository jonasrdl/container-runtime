package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var all bool

var deleteCmd = &cobra.Command{
	Use:   "delete <containerID>",
	Short: "delete a container",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  deleteContainer,
}

func deleteContainer(_ *cobra.Command, args []string) error {
	if all {
		fmt.Print("Are you sure you want to delete all containers? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			return fmt.Errorf("operation cancelled")
		}

		fmt.Println("Deleting all containers...")

		containerDir, err := os.ReadDir(rootPath)
		if err != nil {
			return fmt.Errorf("error reading root path: %v", err)
		}
		for _, container := range containerDir {
			if container.IsDir() {
				containerID := container.Name()
				containerPath := fmt.Sprintf("%s/%s", rootPath, containerID)

				if err := os.RemoveAll(containerPath); err != nil {
					return fmt.Errorf("error deleting container: %v", err)
				}
				fmt.Println("Deleted container: ", containerID)
			}
		}
		fmt.Println("All containers deleted")
	} else {
		if len(args) != 1 {
			return fmt.Errorf("error: container ID not provided")
		}

		containerID := args[0]
		containerPath := fmt.Sprintf("%s/%s", rootPath, containerID)

		if err := os.RemoveAll(containerPath); err != nil {
			return fmt.Errorf("error deleting container: %v", err)
		}
	}
	return nil
}

func init() {
	deleteCmd.Flags().BoolVarP(&all, "all", "a", false, "delete all containers")
	rootCmd.AddCommand(deleteCmd)
}
