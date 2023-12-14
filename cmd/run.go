package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <image> [command]",
	Short: "run a command in a new container",
	Args:  cobra.MinimumNArgs(1),
	RunE:  run,
}

func run(_ *cobra.Command, args []string) error {
	image := args[0]
	var command []string

	if len(args) >= 2 {
		command = args[1:]
	} else {
		defaultCmd, err := getDefaultCommand(image)
		if err != nil {
			return fmt.Errorf("error getting default command: %v", err)
		}
		command = defaultCmd
	}

	fmt.Println("Running", command)

	execCmd := exec.Command("/proc/self/exe", append([]string{"child", image}, command...)...)
	execCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		UidMappings:  []syscall.SysProcIDMap{{HostID: os.Getuid(), Size: 1}},
		GidMappings:  []syscall.SysProcIDMap{{HostID: os.Getgid(), Size: 1}},
	}
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	// Start the command and wait for it to finish
	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("error starting child process: %v", err)
	}
	return nil
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
