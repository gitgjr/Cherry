package video

import (
	"fmt"
	"regexp"
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

func ExtractSerialNumber(filename string) (int, error) {
	re := regexp.MustCompile(`left(\d+)\.ts$`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		return strconv.Atoi(matches[1]) // matches[1] contains the first captured group
	}
	return 0, fmt.Errorf("no serial number found in filename")
}
