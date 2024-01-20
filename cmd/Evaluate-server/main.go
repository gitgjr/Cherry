package main

import (
	"main/utils"
	"main/video"
	"main/zlog"

	"go.uber.org/zap"
)

func main() {

	rootPath := utils.RootPath()
	dataPath := rootPath + "/data"
	// resourcePath := dataPath + "resource"
	serverPath := dataPath + "serverWork"
	durationTime := 5

	zlog.Info("Start time")
	err := video.Mp4toHLS(serverPath+"/left", durationTime)
	if err != nil {
		zlog.Error("mp4 to hls error", zap.Error(err))
	}

	err = video.Mp4toHLS(serverPath+"/right", durationTime)
	if err != nil {
		zlog.Error("mp4 to hls error", zap.Error(err))
	}

	tsFileList, err := utils.FindFiles(serverPath, "", ".ts")
	if err != nil {
		zlog.Error("find ts file error", zap.Error(err))
	}

	if len(tsFileList)%2 != 0 {
		zlog.Error("the number of ts files is not odd")
	}

	err = video.MergeTSFile()

	zlog.Info("End time")
}
