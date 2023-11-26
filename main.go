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
		//w.Update() //test good

		newTransmitTask := mr.TransmitTaskSet{}
		for k, _ := range w.TaskList {
			newTransmitTask["localhost:1116"] = append(newTransmitTask["localhost:1116"], k)
		}
		w.Transmit(newTransmitTask)

		w.Run()

	case "c":
		c := mr.NewCoordinator()
		c.Router()
		c.Run()
	case "w2":
		w2 := mr.NewWorker()
		w2.Port = 1116
		w2.Regester()
		w2.Run()

	}

	//c := mr.Coordinator{}
	//c.Router()
	//c.Run()
}

//TODO: If online test:1. Change workerAddr to public IP 2. Change coordinatorAddr to public IP
