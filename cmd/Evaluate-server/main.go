package main

import (
	"fmt"
	"main/utils"
	"main/video"
	"main/zlog"
	"os"
	"sync"
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
	switch arg1 {
	case "4k":
		leftFile = "4k30left"
		rightFile = "4k30right"
	case "1080":
		leftFile = "1080left"
		rightFile = "1080right"
	default:
	}

	startTime := time.Now()
	if arg2 != "GPU" {
		convertCPU(leftFile, rightFile, serverPath, durationTime)
		// convertCPU_2(leftFile, rightFile, serverPath, durationTime)
		mergeCPU(leftFile, rightFile, serverPath, durationTime)
	}

	execTime := time.Since(startTime)
	fmt.Println("exec time is", execTime.Seconds())
}

func convertCPU(leftFile, rightFile, serverPath string, durationTime int) {
	var wg sync.WaitGroup
	func() {
		wg.Add(1)
		err := video.Mp4toHLS(leftFile, durationTime, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
		wg.Done()

	}()

	func() {
		wg.Add(1)
		err := video.Mp4toHLS(rightFile, durationTime, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
		wg.Done()
	}()
	wg.Wait()
}

func convertCPU_2(leftFile, rightFile, serverPath string, durationTime int) {
	var wg sync.WaitGroup
	func() {
		wg.Add(1)
		err := video.Mp4toHLS_2(leftFile, durationTime, 30, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
		wg.Done()

	}()

	func() {
		wg.Add(1)
		err := video.Mp4toHLS_2(rightFile, durationTime, 30, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
		wg.Done()
	}()
	wg.Wait()
}

func mergeCPU(leftFile, rightFile, serverPath string, durationTime int) {
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
}
