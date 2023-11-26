package mr

import (
	"fmt"
	"io"
	"log"
	"main/utils"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var DataPath = "./data"
var TempPath = "./temp"
var WorkerPort = 1115 //Rename this to WorkerPort
var CoordinatorAddr = "localhost:8080"

type Worker struct {
	WorkerID string
	Addr     string
	TaskList Task
	Port     int
}

func newID() string {
	//Generate a unique 8-digit ID use rand function
	id := rand.Intn(100000000)
	return strconv.Itoa(id)
}

func (w *Worker) Run() {
	http.ListenAndServe(":"+strconv.Itoa(w.Port), nil)
}

func (w *Worker) Router() {
	http.HandleFunc("/", w.DefaultHandler)
	http.HandleFunc("/transmitOrder", w.TransmitOrderHandler)
	//http.HandleFunc("/transmit", w.TransmitHandler)

}

func NewWorker() *Worker {
	w := Worker{}
	w.WorkerID = newID()
	addr, err := utils.GetOutBoundIP() //For online test
	if err != nil {
		panic(err)
	}
	w.Addr = addr + ":" + strconv.Itoa(WorkerPort)
	w.TaskList = make(Task)
	w.Port = WorkerPort
	return &w
}

func (w *Worker) AddMapTask() {
	fileList, err := utils.FindFiles(DataPath, "segment", ".ts")
	if err != nil {
		panic(err)
	}
	for _, file := range fileList {
		fMeta := FileMeta{}
		fMeta.FileName = file
		fMeta.Location = DataPath + "/" + file
		fMeta.FileID, err = utils.GetFileHash(fMeta.Location)
		if err != nil {
			panic(err)
		}
		fMeta.FileSize, err = utils.FileSize(fMeta.Location)
		if err != nil {
			panic(err)
		}

		fMeta.UploadTime = time.Now().Format("2006-01-02 15:04:05")
		w.TaskList[fMeta.FileID] = fMeta
	}

}

func (w *Worker) Regester() {
	//res, err := http.NewRequest(http.MethodGet, CoordinatorAddr+"/register", nil)

	res, err := SendPostRequest(w, "http://"+CoordinatorAddr+"/register")
	if err != nil {
		panic(err)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
}

func (w *Worker) Update() {

	res, err := SendPutRequest(w, "http://"+CoordinatorAddr+"/update")

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
}

// SendTask :send a single file to target worker by reading the whole file in memory
func (w *Worker) sendTask(taskID utils.HashValue, targetAddr string) (*http.Response, error) {
	file, err := os.Open(w.TaskList[taskID].Location)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	newTransmitTask := SingleTransmitTask{
		TaskID: taskID,
		FMeta:  w.TaskList[taskID],
		FData:  data,
	}
	res, err := SendFileRequest(newTransmitTask, "http://"+targetAddr+"/send") //send a single file
	return res, err
}

// Transmit :transmit the task to target workers
// Two methods to transmit:
// 1. Read and send a file one by one
// 2. Use Multipart send some files which divided into chunks in a same time
func (w *Worker) Transmit(tasks TransmitTaskSet) {
	fmt.Println("Transmitting")
	for workerAddr := range tasks {
		files, _ := tasks[workerAddr]
		for _, file := range files {
			_, ok := w.TaskList[file]
			if ok == false {
				fmt.Println("File not exist, send update to coordinator")
				w.Update()
				return
			}
			//Send the file to the worker
			res, err := w.sendTask(file, workerAddr)
			if err != nil {
				log.Fatal(err)
				return
			} else {
				//Write file
				fmt.Println(res)
				continue
			}
		}

	}
}

func (w *Worker) checksumHash(filePath string, fileID utils.HashValue) (bool, error) {
	fileHash, err := utils.GetFileHash(filePath)
	if err != nil {
		return false, err
	}
	if fileHash == fileID {
		return true, nil
	} else {
		return false, nil
	}
}

func (w *Worker) RunServer() {
	http.ListenAndServe(":"+strconv.Itoa(WorkerPort), nil)
}
