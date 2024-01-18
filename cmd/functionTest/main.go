package main

import (
	"main/utils"
	"main/video"
	"os"
)

func main() {

	rootPath := utils.RootPath()
	dataPath := rootPath + "/data/"

	arg1 := os.Args[1]
	switch arg1 {
	case "merge": //good
		videoList := []string{
			dataPath + "/mp4/left_output2.ts",
			dataPath + "/mp4/output2.ts",
		}
		outputFile := dataPath + "m_output2.ts"
		video.StackChunks("vstack", videoList, outputFile)
		// case "HLS": //good
		// 	video.ConvertToHLS(dataPath+"right_added.mp4", dataPath+"right_out.m3u8", 5)
		// case "addKeyFrame": //good
		// 	video.AddKeyFrame(dataPath+"right.mp4", dataPath+"added_right.mp4", 5)
		// case "startTime":
		// 	st, err := video.GetVideoStartTime(dataPath + "right_out0.ts")
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	fmt.Print(st)
		// case "reset": //good
		// 	st, err := video.GetVideoStartTime(dataPath + "right_out0.ts")
		// 	if err != nil {
		// 		fmt.Print(err)
		// 	}
		// 	fmt.Print(st)
		// 	err = video.ResetTimeStamp(dataPath+"m_output2.ts", dataPath+"new_output2.ts", 2, float64(5), st)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
	}

}
