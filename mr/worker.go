package mr

import (
	"fmt"
	"io"
	"log"
	"main/hash"
	"main/httpRequest"
	"main/meta"
	"main/utils"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var rootPath = utils.RootPath()
var DataPath = rootPath + "/data"
var TempPath = rootPath + "/temp"
var LocalAddr = "localhost"
var WorkerPort = "1115 " //Rename this to WorkerPort

var WorkerAddr = LocalAddr
var CoordinatorAddr = LocalAddr
var CoordinatorPort = "8080"

// addr, err := utils.GetOutBoundIP() //For online test
// if err != nil {
// 	panic(err)
// }

type Worker struct {
	WorkerID string
	Addr     string
	Port     string
	TaskList Task
	mutex    sync.Mutex
	//only for coordinator
	State      string //online , offline
	LastOnline time.Time
}

func newID() string {
	//Generate a unique 8-digit ID use rand function
	id := rand.Intn(100000000)
	return strconv.Itoa(id)
}

func (w *Worker) Run() {
	http.ListenAndServe(":"+w.Port, nil)
}

func (w *Worker) Router() {
	http.HandleFunc("/", w.DefaultHandler)
	http.HandleFunc("/transmitOrder", w.TransmitOrderHandler)
	http.HandleFunc("/send", w.TransmitHandler)
	http.HandleFunc("/check", w.CheckHandler)
	http.HandleFunc("/checkState", w.CheckStateHandler)
	//http.HandleFunc("/transmit", w.TransmitHandler)

}

func NewWorker() *Worker {
	w := Worker{}
	w.WorkerID = newID()
	w.Addr = WorkerAddr
	w.Port = WorkerPort
	w.TaskList = make(Task)
	return &w
}

// if filename=="",add all files,else add the file
func (w *Worker) AddMapTask(fileName []string) {
	if len(fileName) == 0 {
		fileList, err := utils.FindFiles(DataPath, "segment", ".ts")
		if err != nil {
			panic(err)
		}
		for _, file := range fileList {
			fMeta, err := meta.FileToFileMeta(file, DataPath+"/"+file)
			if err != nil {
				panic(err)
			}

			w.TaskList[fMeta.FileID] = *fMeta
		}
	} else {
		for _, file := range fileName {
			fMeta, err := meta.FileToFileMeta(file, DataPath+"/"+file)
			if err != nil {
				panic(err)
			}
			w.TaskList[fMeta.FileID] = *fMeta
		}
	}

}

func (w *Worker) readResponse(res *http.Response) {
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
}

func (w *Worker) checkTask(taskID hash.HashValue) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	_, ok := w.TaskList[taskID]
	return ok
}

//--------C/S--------

func (w *Worker) Register() {
	//res, err := http.NewRequest(http.MethodGet, CoordinatorAddr+"/register", nil)

	res, err := httpRequest.SendPostRequest(w, utils.SpliceUrl(CoordinatorAddr, CoordinatorPort, "register"))
	if err != nil {
		panic(err)
	}
	w.readResponse(res)

}

func (w *Worker) Update() {
	res, err := httpRequest.SendPutRequest(w, utils.SpliceUrl(CoordinatorAddr, CoordinatorPort, "update"))
	if err != nil {
		panic(err)
	}
	w.readResponse(res)
}

// CallReduce Since map phrase is auto,it is call for Reduce
func (w *Worker) CallReduce() {
	fmt.Println("Call Reduce")
	res, err := httpRequest.SendPostRequest(w, utils.SpliceUrl(CoordinatorAddr, CoordinatorPort, "callReduce"))
	if err != nil {
		panic(err)
	}
	w.readResponse(res)
}

//--------p2p--------

func (w *Worker) CheckP2PConnect(targetAddr string) {
	res, err := http.Get(targetAddr + "/check")
	if err != nil {
		panic(err)
	}
	w.readResponse(res)
}

// sendTask send a single file to target worker by reading the whole file in memory
func (w *Worker) sendTask(taskID hash.HashValue, targetAddr string) (*http.Response, error) {
	file, err := os.Open(w.TaskList[taskID].Location)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	newTransmitTask := singleTransmitTask{
		TaskID: taskID,
		FMeta:  w.TaskList[taskID],
		FData:  data,
	}
	res, err := httpRequest.SendFileRequest(newTransmitTask, "http://"+targetAddr+"/send") //send a single file
	return res, err
}

// Transmit :transmit the task to target workers
// Two methods to transmit:
// 1. Read and send a file one by one
// 2. Use Multipart send some files which divided into chunks in a same time
func (w *Worker) Transmit(tasks TransmitTask) {
	fmt.Println("Transmitting")
	var wg sync.WaitGroup
	for workerAddr := range tasks {
		files, _ := tasks[workerAddr]

		//create treads for every target
		wg.Add(1)
		go func(targetAddr string, taskID []hash.HashValue) {
			defer wg.Done()
			for _, file := range taskID {
				ok := w.checkTask(file)
				// if
				if !ok {
					fmt.Println("File not exist, send update to coordinator")
					w.Update()
					return
				} else {
					w.mutex.Lock()
					defer w.mutex.Unlock()
					res, err := w.sendTask(file, targetAddr)
					if err != nil {
						log.Fatal(err)
						return
					} else {
						//Write file
						w.readResponse(res)
						continue
					}
				}
			}
		}(workerAddr, files)

	}
	wg.Wait()
}

func (w *Worker) checksumHash(filePath string, fileID hash.HashValue) (bool, error) {
	fileHash, err := hash.GetFileHash(filePath)
	if err != nil {
		return false, err
	}
	if fileHash == fileID {
		return true, nil
	} else {
		return false, nil
	}
}

// func (w *Worker) RunServer() {
// 	http.ListenAndServe(":"+strconv.Itoa(WorkerPort), nil)
// }
