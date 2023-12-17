package main

import (
	"main/nc"
	"main/utils"
	"main/video"
	"os"
)

// TODO:1.Video Merge 2.Implement DHT
func main() {

	rootPath := utils.RootPath()
	dataPath := rootPath + "/data/"

	arg1 := os.Args[1]
	switch arg1 {
	case "nc":
		nc.NcExample()
	case "merge":
		videoList := []string{
			dataPath + "left.mp4",
			dataPath + "right.mp4",
		}
		outputFile := dataPath + "Moutput.mp4"
		video.MergeMP4s(videoList, outputFile)
	case "HLS":
		video.ConvertToHLS(dataPath+"right.mp4", dataPath+"right_out", "2")
	}

}
