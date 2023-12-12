package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"syscall"
)

var runCmd = &cobra.Command{
	Use:   "run [image] [command]",
	Short: "run a command in a new container",
	Args:  cobra.ExactArgs(2),
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	image := args[0]
	command := args[1:]

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

func init() {
	rootCmd.AddCommand(runCmd)
}
