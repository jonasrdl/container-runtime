package cmd

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var childCmd = &cobra.Command{
	Use:    "child",
	Short:  "child is the child process that runs the container",
	Args:   cobra.MinimumNArgs(1),
	Hidden: true,
	RunE:   child,
}

var defaultCommandFlag string

func child(_ *cobra.Command, args []string) error {
	image := args[0]

	// If a command is provided explicitly, use it; otherwise, check for the default command
	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	} else if defaultCommandFlag != "" {
		// If defaultCommandFlag is provided, use it as the command
		cmdArgs = []string{defaultCommandFlag}
	} else {
		// If no explicit command and no default command, use the default command from the image
		defaultCommand, err := getDefaultCommand(image)
		if err != nil {
			return fmt.Errorf("error getting default command: %v", err)
		}
		cmdArgs = defaultCommand
	}

	containerID := generateContainerID()

	// Create a temporary directory to extract the image contents
	tempDir := fmt.Sprintf("/var/lib/container-runtime/%s-tempfs", containerID)

	// Defer the cleanup function to remove the temp directory on exit
	defer cleanupTempDir(tempDir)

	if err := os.MkdirAll(tempDir, 0o770); err != nil {
		return fmt.Errorf("error creating temp directory: %v", err)
	}

	if err := exec.Command("tar", "xvf", "assets/"+image+".tar.gz", "-C", tempDir, "--no-same-owner").Run(); err != nil {
		return fmt.Errorf("error extracting image: %v", err)
	}

	newRootFolder := fmt.Sprintf("/var/lib/container-runtime/%s", containerID)

	// move the temp folder to the root filesystem /var/lib/container-runtime/<containerID>
	if err := os.Rename(tempDir, newRootFolder); err != nil {
		return fmt.Errorf("error moving temp directory: %v", err)
	}

	if err := syscall.Sethostname([]byte(containerID)); err != nil {
		return fmt.Errorf("error setting hostname: %v", err)
	}

	if err := syscall.Chroot(newRootFolder); err != nil {
		return fmt.Errorf("error changing root filesystem: %v", err)
	}

	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("error changing directory: %v", err)
	}

	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("error mounting proc: %v", err)
	}

	fmt.Printf("Running command: %v\n", cmdArgs)
	execCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}

// generateContainerID generates a unique ID for the container
func generateContainerID() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(time.Now().Unix(), 10))))
}

// cleanupTempDir removes the temporary directory
func cleanupTempDir(tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		fmt.Printf("Error cleaning up temp directory %s: %v\n", tempDir, err)
	}
}

func init() {
	childCmd.Flags().StringVar(&defaultCommandFlag, "default-command", "", "default command to run in the container")
	rootCmd.AddCommand(childCmd)
}
