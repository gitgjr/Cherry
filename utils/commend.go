package utils

import (
	"os"
	"os/exec"
	"strings"
)

// ArgsFromCommend Extracting parameters from a command, return the list of args without commend
func ArgsFromCommend(commendString string) []string {
	args := strings.Fields(commendString)
	args = args[1:]
	return args
}

// RunCommend Exec commend with commend name and args
func RunCommend(commendName string, args []string) error {
	cmd := exec.Command(commendName, args...)

	// Redirect command output to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
