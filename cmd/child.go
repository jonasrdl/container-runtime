package cmd

import (
	"crypto/sha256"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

var childCmd = &cobra.Command{
	Use:   "child",
	Short: "child is the child process that runs the container",
	Args:  cobra.ExactArgs(2),
	Run:   child,
}

func child(_ *cobra.Command, args []string) {
	image := args[0]
	command := args[1]

	containerID := generateContainerID()

	// Create a temporary directory to extract the image contents
	tempDir := fmt.Sprintf("./%s-tempfs", containerID)

	must(os.MkdirAll(tempDir, 0770))

	must(exec.Command("tar", "xvf", "assets/"+image+".tar.gz", "-C", tempDir).Run())

	newRootFolder := fmt.Sprintf("/var/lib/container-runtime/%s", containerID)

	// move the temp folder to the root filesystem /var/lib/container-runtime/<containerID>
	must(os.Rename(tempDir, newRootFolder))

	must(syscall.Sethostname([]byte(containerID)))
	must(syscall.Chroot(newRootFolder))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	// create a /dev/null file in the container
	_, err := os.Create("/dev/null")
	if err != nil {
		fmt.Println("ERROR creating /dev/null:", err)
		os.Exit(1)
	}

	fmt.Printf("Running command: %v\n", command)
	execCmd := exec.Command(command)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// generateContainerID generates a unique ID for the container
func generateContainerID() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(time.Now().Unix(), 10))))
}

func init() {
	rootCmd.AddCommand(childCmd)
}
