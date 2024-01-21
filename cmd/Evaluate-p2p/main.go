package main

import (
	"main/utils"
	"main/video"
	"main/zlog"
	"os"
	"slices"

	"go.uber.org/zap"
)

func main() {
	rootPath := utils.RootPath()
	dataPath := rootPath + "/data"
	// resourcePath := dataPath + "resource"
	serverPath := dataPath + "/" + "serverWork"
	durationTime := 5

	leftFile := "left"
	rightFile := "right"

	arg1 := os.Args[1]
	arg2 := os.Args[2]
	switch arg2 {
	case "4k":
		leftFile = "4kleft"
		rightFile = "4kright"
	case "1080":
		leftFile = "1080left"
		rightFile = "1080right"
	default:
	}

	switch arg1 {
	case "convertLeft":
		err := video.Mp4toHLS(leftFile, durationTime, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
	case "convertRight":
		err := video.Mp4toHLS(rightFile, durationTime, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
	case "Merge":
		tsFileList, err := utils.FindFiles(serverPath, "", ".ts")
		if err != nil {
			zlog.Error("find ts file error", zap.Error(err))
		}
		if len(tsFileList)%2 != 0 {
			zlog.Error("the number of ts files is not odd")
		}

		tsFileNumber := len(tsFileList) / 2

		indexList := []int{}

		for _, e := range tsFileList {
			index, err := video.ExtractSerialNumber(e)
			if err != nil {
				zlog.Error("video index error", zap.Error(err))
			}
			if !slices.Contains(indexList, index) {
				indexList = append(indexList, index)
			}
		}

		for i := 0; i < tsFileNumber; i++ {
			tsFilePair := video.FindTsFileByIndex(tsFileList, i)
			err = video.MergeTSFile(tsFilePair, tsFilePair[0], i, "vstack", durationTime, serverPath)
			if err != nil {
				zlog.Error("merge error", zap.Error(err))
			}
		}
	}

}
