package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type HashValue string

func GetFileHash(filePath string) (HashValue, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return HashValue(hashString), nil
}
