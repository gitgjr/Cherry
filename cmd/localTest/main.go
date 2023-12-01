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
		w.AddMapTask(nil)
		w.Register()
		// w.Update() //test good

		// w.CheckP2PConnect(":1116") //test good

		// counter := 0
		// newTransmitTask := mr.TransmitTaskSet{}
		// for k, _ := range w.TaskList {
		// 	if counter == 0 {
		// 		newTransmitTask["localhost:1116"] = append(newTransmitTask["localhost:1116"], k)
		// 	} else {
		// 		newTransmitTask["localhost:1117"] = append(newTransmitTask["localhost:1117"], k)
		// 	}
		// 	counter++
		// }
		// w.Transmit(newTransmitTask) //test good

		w.Router()
		w.Run()

	case "c":
		c := mr.NewCoordinator()
		c.Router()
		c.Run()
	case "w2": //empty worker for p2p test
		w2 := mr.NewWorker()
		w2.Addr = mr.LocalAddr
		w2.Port = "1116"
		w2.Register()
		w2.Router()
		w2.Run()
	case "w3": //empty worker for p2p test
		w3 := mr.NewWorker()
		w3.Addr = mr.LocalAddr
		w3.Port = "1117"
		w3.Register()
		w3.Router()
		w3.Run()
	case "w4": //worker with part of task,simulate mapped worker
		w4 := mr.NewWorker()
		w4.Addr = mr.LocalAddr
		w4.Port = "1118"
		taskList := []string{"segment000.ts"}
		w4.AddMapTask(taskList)
		w4.Register()
		w4.Router()
		w4.Run()
	case "w5": //worker with part of task,simulate mapped worker
		w5 := mr.NewWorker()
		w5.Addr = mr.LocalAddr
		w5.Port = "1119"
		taskList := []string{"segment001.ts"}
		w5.AddMapTask(taskList)
		w5.Register()
		w5.Router()
		w5.Run()
	case "w6": //empty worker for reduce ,this worker call reduce
		w6 := mr.NewWorker()
		w6.Addr = mr.LocalAddr
		w6.Port = "1120"
		w6.Register()
		w6.Router()
		go w6.CallReduce()
		w6.Run()

	}

}

//TODO: If online test:1. Change workerAddr to public IP 2. Change coordinatorAddr to public IP
