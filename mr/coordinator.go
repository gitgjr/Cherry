package mr

import (
	"fmt"
	"log"
	"main/hash"
	"main/httpRequest"
	"main/utils"
	"net/http"
	"sync"
	"time"
)

type Coordinator struct {
	Workers  map[string]*Worker //[workerID]*Worker
	NWorkers int
	Bucket
	allTask Task //all tasks from registered worker
	// TaskChannel chan Task
	mutex sync.Mutex
}

//TODO: Message service and P2P service

type Bucket map[int]hash.HashValue //WorkerID:Task

func NewCoordinator() *Coordinator {
	c := Coordinator{}
	c.Workers = make(map[string]*Worker)
	c.allTask = make(Task)
	return &c
}

// Run Boot Http server
func (c *Coordinator) Run() {
	http.ListenAndServe(":8080", nil)
}

func (c *Coordinator) Router() {
	http.HandleFunc("/", c.DefaultHandler)
	http.HandleFunc("/register", c.RegisterHandler)
	http.HandleFunc("/update", c.UpdateHandler)
	http.HandleFunc("/callReduce", c.callTransmitHandler)
}

func (c *Coordinator) PrintWorkers() {
	for _, worker := range c.Workers {
		utils.PrintStruct(worker)
	}
}

// ScanAllTask:Scan all registered creator and add all the task they hold into allTask
func (c *Coordinator) ScanAllTask() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, w := range c.Workers {
		err := MergeTasks(c.allTask, w.TaskList)
		if err != nil {
			panic(err)
		}
	}
}

// returnOnlineWorker return list of ID of online worker and list of online worker
func (c *Coordinator) returnOnlineWorker() ([]string, []*Worker) {
	onlineIDList := []string{}
	onlineWorkerList := []*Worker{}
	for id, w := range c.Workers {
		if w.State == "online" {
			onlineIDList = append(onlineIDList, id)
			onlineWorkerList = append(onlineWorkerList, w)
		}
	}
	return onlineIDList, onlineWorkerList
}

// returnNonemptyWorker return list of ID of Nonempty worker and list of Nonempty worker
func (c *Coordinator) returnNonemptyWorker() ([]string, []*Worker) {
	nonemptyIDList := []string{}
	nonemptyWorkerList := []*Worker{}
	for id, w := range c.Workers {
		if len(w.TaskList) > 0 {
			nonemptyIDList = append(nonemptyIDList, id)
			nonemptyWorkerList = append(nonemptyWorkerList, w)
		}
	}
	return nonemptyIDList, nonemptyWorkerList
}

// AssignMapWork :
// M1. Average assign without bandwidth
// M2. Bandwidth average assign
func (c *Coordinator) AssignMapWork() {

}

// AssignReduceWork :
// M1.if worker get this task assign(Random assignment)
// M2.Equally distributed according to the number of workers(Average assignment)
// M3.Assign based on worker connections on the basis of 2(RTT assignment)
func (c *Coordinator) AssignReduceTask() ReduceTaskSet {
	c.ScanAllTask()
	_, onlineWorkerList := c.returnOnlineWorker() //[WorkerID]
	newTransmitTaskSet := c.assignReduceTaskM1(onlineWorkerList)
	return newTransmitTaskSet
}

// CheckWorkers: check and update the state of workers
func (c *Coordinator) CheckWorkers() {
	var wg sync.WaitGroup
	for workerID := range c.Workers {
		wg.Add(1)
		go func(wid string) {
			defer wg.Done()
			err := c.sendCheckAndUpdate(wid)
			fmt.Println(err)

		}(workerID)
	}
	wg.Wait()
	// c.PrintWorkers()
}

func (c *Coordinator) sendCheckAndUpdate(workerID string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	targetWorker := c.Workers[workerID]
	res, err := http.Get(utils.SpliceUrl(targetWorker.Addr, targetWorker.Port, "checkState"))
	if err != nil {
		c.Workers[workerID].State = "offline"
		return err
	}
	if res.StatusCode != http.StatusOK {
		c.Workers[workerID].State = "offline"
		return nil
	}
	c.Workers[workerID].State = "online"
	c.Workers[workerID].LastOnline = time.Now()
	return nil
}

// transmit: give a worker a command to transmit receivers:WorkerID of receivers,
// transmitTaskID:TaskID of tasks to be transmitted
func (c *Coordinator) transmit(sender *Worker, tTask TransmitTask) {
	//check if task exist
	for _, taskIDList := range tTask {
		for _, taskID := range taskIDList {
			_, ok := sender.TaskList[taskID]
			if ok == false {
				panic("task not found" + taskID)
			}
		}

	}
	res, err := httpRequest.SendPostRequest(tTask, utils.SpliceUrl(sender.Addr, sender.Port, "transmitOrder"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res) //TODO: Change here
}

func (c *Coordinator) DivideBucket() {

}

// scanForFiles:Add files into MTask list via prefix and suffix
