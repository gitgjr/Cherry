package mr

import (
	"errors"
	"main/hash"
	"main/meta"
)

type Task map[hash.HashValue]meta.FileMeta //taskID : fileMeta

// WorkerAddr:[]TaskID, let one worker send its files to multiple  workers
// give a sender a list :receiver.addr:[]TaskID
// worker1.addr:[task1,task2],worker2.addr:[task3,task4]
type TransmitTask map[string][]hash.HashValue

// WorkerID:[]TaskID,let multiple  workers to send their files one worker
// give multiple sender a list :sender.addr:[]TaskID ,set means it needs to be divided into several TransmitTask
// worker1.addr:[task1,task2],worker2.addr:[task3,task4]
type ReduceTaskSet map[string][]hash.HashValue

type SingleTransmitTask struct {
	TaskID hash.HashValue
	FMeta  meta.FileMeta
	FData  []byte
}

// MakeTransmitTask :For all tasks, scan all workers
// and let the task be transmitted by that worker if that worker has that task.
// The crudest way to assign
func MakeTransmitTask(receivers []*Worker, transmitTaskID []hash.HashValue) TransmitTask {
	t := make(TransmitTask)
	for _, receiver := range receivers {
		for _, taskID := range transmitTaskID {
			_, receiverOk := receiver.TaskList[taskID]
			_, added := t[string(taskID)]
			if !added && receiverOk {
				t[receiver.Addr] = append(t[receiver.Addr], taskID)
			}
		}
	}
	return t
}

// MergeTasks: Merge m2 to m1,return error when they have same keys
func MergeTasks(m1 Task, m2 Task) error {
	for k, v := range m2 {
		if _, ok := m1[k]; ok {
			return errors.New("hash conflict, two hashmap have same key")
		}
		m1[k] = v
	}
	return nil
}
