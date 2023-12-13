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
	Args:   cobra.MinimumNArgs(2),
	Hidden: true,
	RunE:   child,
}

func child(_ *cobra.Command, args []string) error {
	image := args[0]
	command := args[1]

	containerID := generateContainerID()

	// Create a temporary directory to extract the image contents
	tempDir := fmt.Sprintf("./%s-tempfs", containerID)

	if err := os.MkdirAll(tempDir, 0o770); err != nil {
		return fmt.Errorf("error creating temp directory: %v", err)
	}

	if err := exec.Command("tar", "xvf", "assets/"+image+".tar.gz", "-C", tempDir).Run(); err != nil {
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

	// create a /dev/null file in the container
	if _, err := os.Create("/dev/null"); err != nil {
		return fmt.Errorf("error creating /dev/null: %v", err)
	}

	fmt.Printf("Running command: %v\n", command)
	execCmd := exec.Command(command)
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

func init() {
	rootCmd.AddCommand(childCmd)
}
