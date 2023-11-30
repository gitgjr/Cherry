package mr

import (
	"encoding/json"
	"fmt"
	"io"
	"main/utils"
	"net/http"
)

func (c *Coordinator) DefaultHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "This is the default handler")
}

func UploadHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Println("i haven`t design get")
	} else if req.Method == http.MethodPost {
		//upload
	}
}

func (c *Coordinator) RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Println("Only post method is allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		var newWorker Worker
		err := json.NewDecoder(req.Body).Decode(&newWorker)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c.Workers[newWorker.WorkerID] = &newWorker
		//fmt.Println(newWorker)
		io.WriteString(w, "register worker success")
	}
}

func (c *Coordinator) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		fmt.Println("Only put method is allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		var updateWorker Worker
		err := json.NewDecoder(req.Body).Decode(&updateWorker)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c.Workers[updateWorker.WorkerID] = &updateWorker
		//fmt.Println(updateWorker)
		io.WriteString(w, "update worker success")
	}
}

// Reduce phrase
// 1.Check online worker
// 2.Assign task
// 3.Send task
func (c *Coordinator) callTransmitHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Println("Only Post method is allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		var requestWorker Worker
		err := json.NewDecoder(req.Body).Decode(&requestWorker)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c.CheckWorkers() //test good
		taskSet := c.AssignReduceTask()
		for senderID, taskList := range taskSet {
			//create transmit task
			go func(sID string, task []utils.HashValue) {
				newTransmitTask := make(TransmitTask)
				newTransmitTask[requestWorker.WorkerID] = task
				//send order
				c.transmit(c.Workers[sID], newTransmitTask)
			}(senderID, taskList)

		}

		// for k,v:=range c.Workers{
		// 	c.transmit()
		// }

	}
}
