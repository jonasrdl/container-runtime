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
	Run:   deleteContainer,
}

func deleteContainer(_ *cobra.Command, args []string) {
	if all {
		fmt.Print("Are you sure you want to delete all containers? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		if strings.TrimSpace(input) != "y" {
			fmt.Println("Operation cancelled")
			os.Exit(0)
		}

		fmt.Println("Deleting all containers...")

		containerDir, err := os.ReadDir(rootPath)
		if err != nil {
			fmt.Println("Error reading root path:", err)
			os.Exit(1)
		}
		for _, container := range containerDir {
			if container.IsDir() {
				containerID := container.Name()
				containerPath := fmt.Sprintf("%s/%s", rootPath, containerID)

				if err := os.RemoveAll(containerPath); err != nil {
					fmt.Println("Error deleting container:", err)
					os.Exit(1)
				}
				fmt.Println("Deleted container: ", containerID)
			}
		}
		fmt.Println("All containers deleted")
	} else {
		if len(args) != 1 {
			fmt.Println("Error: container ID not provided")
			os.Exit(1)
		}

		containerID := args[0]
		containerPath := fmt.Sprintf("%s/%s", rootPath, containerID)

		if err := os.RemoveAll(containerPath); err != nil {
			fmt.Println("Error deleting container:", err)
			os.Exit(1)
		}
	}
}

func init() {
	deleteCmd.Flags().BoolVarP(&all, "all", "a", false, "delete all containers")
	rootCmd.AddCommand(deleteCmd)
}
