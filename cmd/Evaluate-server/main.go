package main

import (
	"fmt"
	"main/utils"
	"main/video"
	"main/zlog"
	"os"
	"time"

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
	switch arg1 {
	case "4k":
		leftFile = "4kleft"
		rightFile = "4kright"
	case "1080":
		leftFile = "1080left"
		rightFile = "1080right"
	}

	startTime := time.Now()
	err := video.Mp4toHLS(leftFile, durationTime, serverPath)
	if err != nil {
		zlog.Error("mp4 to hls error", zap.Error(err))
	}

	err = video.Mp4toHLS(rightFile, durationTime, serverPath)
	if err != nil {
		fmt.Println("mp4 to hls error", zap.Error(err))
	}

	tsFileList, err := utils.FindFiles(serverPath, "", ".ts")
	if err != nil {
		zlog.Error("find ts file error", zap.Error(err))
	}

	if len(tsFileList)%2 != 0 {
		zlog.Error("the number of ts files is not odd")
	}

	tsFileNumber := len(tsFileList) / 2

	for i := 0; i < tsFileNumber; i++ {
		tsFilePair := video.FindTsFileByIndex(tsFileList, i)
		err = video.MergeTSFile(tsFilePair, tsFilePair[0], i, "vstack", durationTime, serverPath)
		if err != nil {
			zlog.Error("merge error", zap.Error(err))
		}
	}
	err = video.NewM3u8(serverPath+"/"+leftFile+".m3u8", serverPath+"/new_"+leftFile+".m3u8")
	if err != nil {
		zlog.Error("generate m3u8 error", zap.Error(err))
	}

	endTime := time.Now()
	execTime := endTime.Sub(startTime)
	fmt.Println("exec time is", execTime.Seconds())
}
