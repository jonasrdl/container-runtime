package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var runCmd = &cobra.Command{
	Use:   "run [image] [command]",
	Short: "run a command in a new container",
	Args:  cobra.RangeArgs(1, 2),
	Run:   run,
}

func run(_ *cobra.Command, args []string) {
	image := args[0]
	var command []string

	if len(args) == 2 {
		command = args[1:]
	} else {
		defaultCmd, err := getDefaultCommand(image)
		if err != nil {
			fmt.Println("Error getting default command:", err)
			os.Exit(1)
		}
		command = defaultCmd
	}

	fmt.Println("Running", command)

	execCmd := exec.Command("/proc/self/exe", append([]string{"child", image}, command...)...)
	execCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	// Start the command and wait for it to finish
	if err := execCmd.Start(); err != nil {
		fmt.Println("ERROR starting child process:", err)
		os.Exit(1)
	}
	if err := execCmd.Wait(); err != nil {
		fmt.Println("ERROR waiting for child process:", err)
		os.Exit(1)
	}
}

func getDefaultCommand(image string) ([]string, error) {
	defaultCmdFile := fmt.Sprintf("assets/%s-cmd", image)
	content, err := os.ReadFile(defaultCmdFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read default command file: %w", err)
	}

	defaultCmd := strings.TrimSpace(string(content))

	return strings.Fields(defaultCmd), nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
