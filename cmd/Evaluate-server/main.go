package main

import (
	"main/utils"
	"main/video"
	"main/zlog"
)

func main() {

	rootPath := utils.RootPath()
	dataPath := rootPath + "/data/"
	resourcePath := rootPath + "/resource/"
	durationTime := 5

	zlog.Info("Start test")
	video.Mp4toHLS(resourcePath+"left.mp4", durationTime)
}
