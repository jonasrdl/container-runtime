package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: container-runtime run [image] [command]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func run() {
	image := os.Args[2]
	command := os.Args[3:]

	fmt.Println("Running", command)

	cmd := exec.Command("/proc/self/exe", append([]string{"child", image}, command...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command and wait for it to finish
	if err := cmd.Start(); err != nil {
		fmt.Println("ERROR starting child process:", err)
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("ERROR waiting for child process:", err)
		os.Exit(1)
	}
}

func child() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: container-runtime run [image] [command]")
		os.Exit(1)
	}

	image := os.Args[2]
	command := os.Args[3:]

	fmt.Println(command)

	// Create a temporary directory to extract the image contents
	tempDir := "./tempfs"
	must(os.MkdirAll(tempDir, 0770))

	must(exec.Command("tar", "xvf", "assets/"+image+".tar.gz", "-C", tempDir).Run())

	newRootFolder := "/var/lib/container-runtime/test"

	// move the temp folder to the root filesystem /var/lib/container-runtime/test/
	must(os.Rename(tempDir, newRootFolder))

	must(syscall.Sethostname([]byte("my-test-container")))
	must(syscall.Chroot(newRootFolder))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	// create a /dev/null file in the container
	os.Create("/dev/null")

	// Run an interactive shell within the container for debugging
	fmt.Printf("Running command: %v\n", command)
	cmd := exec.Command(command[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
