package main

import (
	"main/utils"
	"main/video"
	"main/zlog"
	"os"

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
		leftFile = "4k30left"
		rightFile = "4k30right"
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
		video.Merge_CPU(leftFile, rightFile, serverPath, durationTime)
	}

}
