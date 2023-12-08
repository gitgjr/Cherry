package mr

import (
	"errors"
	"main/hash"
	"main/meta"
)

type task map[hash.HashValue]meta.FileMeta //taskID : fileMeta

// transmitTask WorkerAddr:[]TaskID, let one worker send its files to multiple  workers
// give a sender a list :receiver.addr:[]TaskID
// worker1.addr:[task1,task2],worker2.addr:[task3,task4]
type transmitTask map[string][]hash.HashValue

// mapTaskSet WorkerID:transmitTask , a set of transmit task assigned by coordinator
type mapTaskSet map[string]transmitTask

// reduceTaskSet WorkerID:[]TaskID,let multiple  workers to send their files one worker
// give multiple sender a list :sender.addr:[]TaskID ,set means it needs to be divided into several TransmitTask
// worker1.addr:[task1,task2],worker2.addr:[task3,task4]
type reduceTaskSet map[string][]hash.HashValue

// singleTransmitTask the real data struck for transmit , always convert to json
type singleTransmitTask struct {
	TaskID hash.HashValue
	FMeta  meta.FileMeta
	FData  []byte
}

// makeTransmitTask :For all tasks, scan all workers
// and let the task be transmitted by that worker if that worker has that task.
// The crudest way to assign
func makeTransmitTask(receivers []*Worker, transmitTaskID []hash.HashValue) transmitTask {
	t := make(transmitTask)
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

// mergeTasks: Merge m2 to m1,return error when they have same keys
func mergeTasks(m1 task, m2 task) error {
	for k, v := range m2 {
		if _, ok := m1[k]; ok {
			return errors.New("hash conflict, two hashmap have same key")
		}
		m1[k] = v
	}
	return nil
}

// return the total size of task
func (t task) sumSize() (int, error) {
	if len(t) < 1 {
		return -1, errors.New("Task is empty")
	}

	sum := 0
	for _, meta := range t {
		sum += meta.FileSize
	}
	if sum < 1 {
		return -1, errors.New("sum is zero")
	}
	return sum, nil
}
