package main

import (
	"fmt"
	"main/utils"
	"main/video"
	"main/zlog"
	"time"

	"go.uber.org/zap"
)

func main() {

	rootPath := utils.RootPath()
	dataPath := rootPath + "/data"
	// resourcePath := dataPath + "resource"
	serverPath := dataPath + "/" + "serverWork"
	durationTime := 5

	startTime := time.Now()
	err := video.Mp4toHLS("left", durationTime, serverPath)
	if err != nil {
		zlog.Error("mp4 to hls error", zap.Error(err))
	}

	err = video.Mp4toHLS("right", durationTime, serverPath)
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
	err = video.NewM3u8(serverPath+"/left.m3u8", serverPath+"/new_left.m3u8")
	if err != nil {
		zlog.Error("generate m3u8 error", zap.Error(err))
	}

	endTime := time.Now()
	execTime := endTime.Sub(startTime)
	fmt.Println("exec time is", execTime.Seconds())
}
