package mr

import (
	"main/utils"
)

type FileMeta struct {
	FileID     utils.HashValue
	FileName   string
	FileSize   int
	Location   string //workerID:path
	UploadTime string
}

var fileMates map[utils.HashValue]FileMeta

func init() {
	fileMates = make(map[utils.HashValue]FileMeta)
}

// UpdateFileMeta update info
func UpdateFileMeta(fMeta *FileMeta) {
	fileMates[fMeta.FileID] = *fMeta
}

// GetFileMeta :Get pointer of fileMeta struct by file id
func GetFileMeta(fileID utils.HashValue) *FileMeta {
	result, exist := fileMates[fileID]
	if exist == false {
		return nil
	} else {
		return &result
	}
}
