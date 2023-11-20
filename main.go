package main

import (
	"main/mr"
	"os"
)

func main() {
	//inputFile := "data/RecordRTC-2023109-a2ujrvzyrtl.mp4" // Replace with your input MP4 file
	//outputDirectory := "data"                             // Replace with your desired output directory
	//
	//err := video.ConvertToHLS(inputFile, outputDirectory)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//}

	//w := mr.NewWorker()
	//w.AddMapTask()
	//fmt.Println(w)

	arg1 := os.Args[1]
	switch arg1 {
	case "w":
		w := mr.NewWorker()
		w.AddMapTask()
		w.Regester()

	case "c":
		c := mr.Coordinator{}
		c.Router()
		c.Run()
	}

	//c := mr.Coordinator{}
	//c.Router()
	//c.Run()
}
