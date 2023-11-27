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
		w.Addr = mr.LocalAddr
		w.AddMapTask()
		w.Regester()
		// w.Update() //test good

		// w.CheckP2PConnect(":1116") //test good

		counter := 0
		newTransmitTask := mr.TransmitTaskSet{}
		for k, _ := range w.TaskList {
			if counter == 0 {
				newTransmitTask["localhost:1116"] = append(newTransmitTask["localhost:1116"], k)
			} else {
				newTransmitTask["localhost:1117"] = append(newTransmitTask["localhost:1117"], k)
			}
			counter++
		}
		w.Transmit(newTransmitTask) //test good

		w.Router()
		w.Run()

	case "c":
		c := mr.NewCoordinator()
		c.Router()
		c.Run()
	case "w2":
		w2 := mr.NewWorker()
		w2.Addr = mr.LocalAddr
		w2.Port = 1116
		w2.Regester()
		w2.Router()
		w2.Run()
	case "w3":
		w3 := mr.NewWorker()
		w3.Addr = mr.LocalAddr
		w3.Port = 1117
		w3.Regester()
		w3.Router()
		w3.Run()
	}

}

//TODO: If online test:1. Change workerAddr to public IP 2. Change coordinatorAddr to public IP
