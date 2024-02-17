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

	startTime := time.Now()
	switch arg1 {
	case "convert":
		err := video.Mp4toHLS(leftFile, durationTime, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
		err = video.Mp4toHLS(rightFile, durationTime, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
	case "Merge":
		video.Merge_CPU(leftFile, rightFile, serverPath, durationTime)
	}
	// if arg3 != "GPU" {
	// 	video.Convert_CPU(leftFile, rightFile, serverPath, durationTime)
	// 	// video.Convert_CPU_2(leftFile, rightFile, serverPath, durationTime)
	// 	video.Merge_CPU(leftFile, rightFile, serverPath, durationTime)
	// }

	execTime := time.Since(startTime)
	fmt.Println("exec time is", execTime.Seconds())
}
