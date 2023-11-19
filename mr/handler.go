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

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Println("i haven`t design get")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else if req.Method == http.MethodPost {
		io.WriteString(w, "register worker success")
		var newWorker Worker
		err := json.NewDecoder(req.Body).Decode(&newWorker)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(newWorker)

	}

}

func DefaultHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello")
}
