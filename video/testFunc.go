package video

import (
	"main/utils"
	"main/zlog"
	"sync"

	"go.uber.org/zap"
)

func Convert_CPU(leftFile, rightFile, serverPath string, durationTime int) {
	var wg sync.WaitGroup
	func() {
		wg.Add(1)
		err := Mp4toHLS(leftFile, durationTime, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
		wg.Done()

	}()

	func() {
		wg.Add(1)
		err := Mp4toHLS(rightFile, durationTime, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
		wg.Done()
	}()
	wg.Wait()
}

func Convert_CPU_2(leftFile, rightFile, serverPath string, durationTime int) {
	var wg sync.WaitGroup
	func() {
		wg.Add(1)
		err := Mp4toHLS_2(leftFile, durationTime, 30, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
		wg.Done()

	}()

	func() {
		wg.Add(1)
		err := Mp4toHLS_2(rightFile, durationTime, 30, serverPath)
		if err != nil {
			zlog.Error("mp4 to hls error", zap.Error(err))
		}
		wg.Done()
	}()
	wg.Wait()
}

func Merge_CPU(leftFile, rightFile, serverPath string, durationTime int) {
	tsFileList, err := utils.FindFiles(serverPath, "", ".ts")
	if err != nil {
		zlog.Error("find ts file error", zap.Error(err))
	}
	if len(tsFileList)%2 != 0 {
		zlog.Error("the number of ts files is not odd")
	}

	tsFileNumber := len(tsFileList) / 2

	for i := 0; i < tsFileNumber; i++ {
		tsFilePair := FindTsFileByIndex(tsFileList, i)
		err = MergeTSFile(tsFilePair, tsFilePair[0], i, "vstack", durationTime, serverPath)
		if err != nil {
			zlog.Error("merge error", zap.Error(err))
		}
	}
	err = NewM3u8(serverPath+"/"+leftFile+".m3u8", serverPath+"/new_"+leftFile+".m3u8")
	if err != nil {
		zlog.Error("generate m3u8 error", zap.Error(err))
	}
}

func Merge_CPU_no_m3u8(leftFile, rightFile, serverPath string, durationTime int) {
	tsFileList, err := utils.FindFiles(serverPath, "", ".ts")
	if err != nil {
		zlog.Error("find ts file error", zap.Error(err))
	}
	if len(tsFileList)%2 != 0 {
		zlog.Error("the number of ts files is not odd")
	}

	tsFileNumber := len(tsFileList) / 2

	for i := 0; i < tsFileNumber; i++ {
		tsFilePair := FindTsFileByIndex(tsFileList, i)
		err = MergeTSFile(tsFilePair, tsFilePair[0], i, "vstack", durationTime, serverPath)
		if err != nil {
			zlog.Error("merge error", zap.Error(err))
		}
	}
}
