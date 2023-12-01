package meta

import (
	"main/hash"
	"main/utils"
	"time"
)

type FileMeta struct {
	FileID     hash.HashValue
	FileName   string
	FileSize   int
	Location   string //workerID:path
	UploadTime string
}

var fileMates map[hash.HashValue]FileMeta

func init() {
	fileMates = make(map[hash.HashValue]FileMeta)
}

// UpdateFileMeta update info
func UpdateFileMeta(fMeta *FileMeta) {
	fileMates[fMeta.FileID] = *fMeta
}

// GetFileMeta :Get pointer of fileMeta struct by file id
func GetFileMeta(fileID hash.HashValue) *FileMeta {
	result, exist := fileMates[fileID]
	if exist == false {
		return nil
	} else {
		return &result
	}
}

func FileToFileMeta(fileName, filePath string) (*FileMeta, error) {
	fMeta := FileMeta{}
	fMeta.FileName = fileName
	fMeta.Location = filePath
	var err error
	fMeta.FileID, err = hash.GetFileHash(fMeta.Location)
	if err != nil {
		return nil, err
	}
	fMeta.FileSize, err = utils.FileSize(fMeta.Location)
	if err != nil {
		return nil, err
	}

	fMeta.UploadTime = time.Now().Format("2006-01-02 15:04:05")
	return &fMeta, nil
}
