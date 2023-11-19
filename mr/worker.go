package mr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/utils"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var workPath = "./data"
var port = 1115
var coordinatorAddr = "localhost:8080"

type Worker struct {
	WorkerID string
	Addr     string
	TaskList Task
}

func newID() string {
	//Generate a unique 8-digit ID use rand function
	id := rand.Intn(100000000)
	return strconv.Itoa(id)
}

func NewWorker() *Worker {
	w := Worker{}
	w.WorkerID = newID()
	addr, err := utils.GetOutBoundIP()
	if err != nil {
		panic(err)
	}
	w.Addr = addr + ":" + strconv.Itoa(port)
	w.TaskList = make(Task)
	return &w
}

func (w *Worker) AddMapTask() {
	fileList, err := utils.FindFiles(workPath, "segment", ".ts")
	if err != nil {
		panic(err)
	}
	for _, file := range fileList {
		fMeta := FileMeta{}
		fMeta.FileName = file
		fMeta.FileID, err = utils.GetFileHash(workPath + "/" + file)
		if err != nil {
			panic(err)
		}
		fMeta.FileSize, err = utils.FileSize(workPath + "/" + file)
		if err != nil {
			panic(err)
		}
		fMeta.Location = w.WorkerID + ":" + file
		fMeta.UploadTime = time.Now().Format("2006-01-02 15:04:05")
		w.TaskList[fMeta.FileID] = fMeta
	}

}

func (w *Worker) Regester() {
	//res, err := http.NewRequest(http.MethodGet, coordinatorAddr+"/register", nil)
	jsonData, err := json.Marshal(w)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	bodyReader := bytes.NewReader(jsonData)
	req, err := http.NewRequest(http.MethodPost, "http://"+coordinatorAddr+"/register", bodyReader)
	if err != nil {
		panic(err)
	}
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
}
