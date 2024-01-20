package video

import (
	"strconv"
	"strings"
)

func FindTsFileByIndex(arr []string, x int) []string {
	suffix := strconv.Itoa(x) + ".ts"
	var result []string
	for _, str := range arr {
		if strings.HasSuffix(str, suffix) {
			result = append(result, str)
		}
	}
	return result
}
