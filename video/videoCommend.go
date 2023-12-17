package video

import "strings"

//argsFromCommend Extracting parameters from a command, return a
func argsFromCommend(commendString string) []string {
	args := strings.Fields(commendString)
	args = args[1:]
	return args
}
