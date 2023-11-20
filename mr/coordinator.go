package mr

import (
	"log"
	"main/utils"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Coordinator struct {
	WorkerList []*Worker
	MapTask    Task
	ReduceTask Task
}

type Task map[utils.HashValue]FileMeta

type Bucket struct {
	BucketID int
	TaskList Task
}

// Boot Http server
func (c *Coordinator) Run() {
	http.ListenAndServe(":8080", nil)
}

func (c *Coordinator) Router() {
	http.HandleFunc("/", DefaultHandler)
	http.HandleFunc("/register", RegisterHandler)

}

func (c *Coordinator) Regester() {

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

func NewCoordinator() *Coordinator {
	c := Coordinator{}
	return &c
}

// scanForFiles:Add files into MTask list via prefix and suffix
