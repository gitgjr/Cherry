package mr

import (
	"errors"
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
	Workers            map[string]*Worker //[workerID]*Worker
	ReplicaCoefficient float32
	allTask            task //all tasks from registered worker.
	// After the map process, all workers actually hold 1.2 to 2 times the tasks
	// TaskChannel chan Task
	mutex sync.Mutex
}

type Bucket map[int]hash.HashValue //WorkerID:Task

func NewCoordinator() *Coordinator {
	c := Coordinator{}
	c.ReplicaCoefficient = 1.2
	c.Workers = make(map[string]*Worker)
	c.allTask = make(task)
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

// addAllTask:Scan all registered creator and add all the task they hold into allTask.
// Should be called after all creator registered or ten mins after coordinator start
func (c *Coordinator) addAllTask() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, w := range c.Workers {
		err := mergeTasks(c.allTask, w.TaskList)
		if err != nil {
			panic(err)
		}
	}
}

// returnOnlineWorker return list of ID of online worker and list of online worker
func (c *Coordinator) returnOnlineWorker() ([]string, []*Worker, error) {
	onlineIDList := []string{}
	onlineWorkerList := []*Worker{}
	for id, w := range c.Workers {
		if w.State == "online" {
			onlineIDList = append(onlineIDList, id)
			onlineWorkerList = append(onlineWorkerList, w)
		}
	}
	if len(onlineIDList) < 1 {
		return nil, nil, errors.New("no online worker exists")
	} else {
		return onlineIDList, onlineWorkerList, nil
	}

}

// returnNonemptyWorker return list of ID of Nonempty worker and list of Nonempty worker
func (c *Coordinator) returnNonemptyWorker() ([]string, []*Worker, error) {

	nonemptyIDList := []string{}
	nonemptyWorkerList := []*Worker{}
	for id, w := range c.Workers {
		if len(w.TaskList) > 0 {
			nonemptyIDList = append(nonemptyIDList, id)
			nonemptyWorkerList = append(nonemptyWorkerList, w)
		}
	}
	if len(nonemptyIDList) < 1 {
		return nil, nil, errors.New("no nonempty worker exists")
	} else {
		return nonemptyIDList, nonemptyWorkerList, nil
	}

}

//---------FREEZE BEGIN---------

// // assignMapTask :
// // M1. Bandwidth average assign
// // M2. Bandwidth RTT assign
// func (c *Coordinator) assignMapTask() (mapTaskSet, error) {
// 	c.addAllTask()
// 	nonemptyWorkerID, _, err := c.returnNonemptyWorker()
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, onlineWorker, err := c.returnOnlineWorker()
// 	if err != nil {
// 		return nil, err
// 	}
// 	newMapTaskSet := c.assignMapTaskM1(nonemptyWorkerID, onlineWorker)
// 	return newMapTaskSet, nil
// }

//---------FREEZE END---------

// assignReduceTask :
// M1.if worker get this task assign(Round-rand assignment)
// M2.Equally distributed according to the number of workers(Average assignment)
// M3.Assign based on worker connections on the basis of 2(RTT assignment)
func (c *Coordinator) assignReduceTask() (reduceTaskSet, error) {
	c.addAllTask()                                     //TODO should remove,just for test
	_, onlineWorkerList, err := c.returnOnlineWorker() //[WorkerID]
	if err != nil {
		return nil, err
	}
	newTransmitTaskSet := c.assignReduceTaskM1(onlineWorkerList)
	return newTransmitTaskSet, nil
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

// sendCheckAndUpdate atomic function
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
func (c *Coordinator) transmit(sender *Worker, tTask transmitTask) {
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
