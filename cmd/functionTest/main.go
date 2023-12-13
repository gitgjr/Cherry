package main

import (
	"main/nc"
	"main/utils"
	"main/video"
	"os"
)

func main() {

	rootPath := utils.RootPath()
	dataPath := rootPath + "/data/"

	arg1 := os.Args[1]
	switch arg1 {
	case "nc":
		nc.NcExample()
	case "merge":
		videoList := []string{
			dataPath + "Mvideo1.mp4",
			dataPath + "Mvideo2.mp4",
		}
		outputFile := dataPath + "output.mp4"
		video.MergeMP4s(videoList, outputFile)
	}

}
