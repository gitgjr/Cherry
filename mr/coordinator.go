package mr

import (
	"fmt"
	"log"
	"main/utils"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

type Coordinator struct {
	Workers  map[string]*Worker //[workerID]*Worker
	NWorkers int
	Bucket
	NMapTask    Task
	ReduceTask  Task
	TaskChannel chan Task
	mutex       sync.Mutex
}

//TODO: Message service and P2P service

type Bucket map[int]Task //WorkerID:Task

func NewCoordinator() *Coordinator {
	c := Coordinator{}
	c.Workers = make(map[string]*Worker)
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
}

func (c *Coordinator) PrintWorkers() {
	for _, worker := range c.Workers {
		fmt.Println(worker)
	}
}

// transmit: give a worker a command to transmit receivers:WorkerID of receivers, transmitTaskID:TaskID of tasks to be transmitted
func (c *Coordinator) transmit(sender *Worker, receivers []*Worker, transmitTaskID []utils.HashValue) {
	for _, taskID := range transmitTaskID {
		_, ok := sender.TaskList[taskID]
		if ok == false {
			log.Fatal("Task not exist")
			return
		}
	}

	NewTransmitTask := MakeTransmitTaskSet(receivers, transmitTaskID)
	res, err := SendPostRequest(NewTransmitTask, "http://"+sender.Addr+":"+strconv.Itoa(WorkerPort)+"/transmit")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res) //TODO: Change here
}

func (c *Coordinator) DivideBucket() {
	// TODO: If need to storage distributively
}

func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// scanForFiles:Add files into MTask list via prefix and suffix
