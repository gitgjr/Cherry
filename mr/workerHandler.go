package mr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (c *Coordinator) DefaultHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "This is the default handler")
}

func (worker *Worker) DefaultHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (worker *Worker) TransmitOrderHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Println("Only post method is allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		var transmitTask TransmitTaskSet
		err := json.NewDecoder(req.Body).Decode(&transmitTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//Send files
	}
}

//workflows:1.download into temp folder
// 2.checksum hash
// 3.add to TaskList
// 4.copy to data folder
//

func (worker *Worker) TransmitHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		fmt.Println("Only post method is allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		var transmitTask SingleTransmitTask
		err := json.NewDecoder(req.Body).Decode(&transmitTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tempFilePath := "temp" + "/" + transmitTask.FMeta.FileName
		file, err := os.Create(tempFilePath)
		defer file.Close()
		if err != nil {
			panic(err)
		}
		_, err = file.Write(transmitTask.FData)
		if err != nil {
			panic(err)
		}
		checksumResult, err := worker.checksumHash(tempFilePath, transmitTask.TaskID)
		if err != nil {
			panic(err)
		}
		if checksumResult == false {
			fmt.Println("checksum failed")
			return
		}
		worker.TaskList[transmitTask.TaskID] = transmitTask.FMeta
		//copy to data folder
		err = os.WriteFile("data"+"/"+"new-"+transmitTask.FMeta.FileName, transmitTask.FData, 0644)
		if err != nil {
			panic(err)
		}
	}
}
