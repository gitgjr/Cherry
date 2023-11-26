package mr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
