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

// ExtractSerialNumbers extracts serial numbers from a slice of filenames.
func ExtractTSPairs(filenames []string) (map[int][]string, error) {
	tsMap := make(map[int][]string)

	re := regexp.MustCompile(`(\d+)\.ts$`)

	for _, filename := range filenames {
		matches := re.FindStringSubmatch(filename)
		if len(matches) > 1 {
			number, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, err
			}
			tsMap[number] = append(tsMap[number], filename)
		}
	}

	for index, files := range tsMap {
		if len(files)%2 != 0 {
			return nil, fmt.Errorf("uneven file count for index %d", index)
		}
	}
	return tsMap, nil
}
