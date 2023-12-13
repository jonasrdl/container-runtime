package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var execCmd = &cobra.Command{
	Use:   "exec <containerID> <command>",
	Short: "exec a command in a running container",
	Args:  cobra.ExactArgs(2),
	Run:   execCommand,
}

func execCommand(_ *cobra.Command, args []string) {
	containerID := args[0]
	command := args[1:]

	rootFolder := fmt.Sprintf("/var/lib/container-runtime/%s", containerID)
	must(os.Chdir(rootFolder))

	execCmd := exec.Command(command[0], command[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(execCmd)
}
