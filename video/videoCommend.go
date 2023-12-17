package video

import (
	"os"
	"os/exec"
	"strings"
)

// argsFromCommend Extracting parameters from a command, return a
func argsFromCommend(commendString string) []string {
	args := strings.Fields(commendString)
	args = args[1:]
	return args
}

func runCommend(commend string, args []string) error {
	cmd := exec.Command("ffmpeg", args...)

	// Redirect command output to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
