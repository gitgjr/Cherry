package utils

import (
	"os"
	"path/filepath"
)

func FileSize(filePath string) (int, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return int(fileInfo.Size()), nil
}

func FindFiles(folderPath, prefix, suffix string) ([]string, error) {
	var matchingFiles []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileName := info.Name()

			if (prefix == "" || (prefix != "" && hasPrefix(fileName, prefix))) &&
				(suffix == "" || (suffix != "" && hasSuffix(fileName, suffix))) {
				matchingFiles = append(matchingFiles, fileName)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matchingFiles, nil
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
