package mr

import (
	"main/hash"
	"math/rand"
	"time"
)

//----- MAP -----

// addReplica randomly add replicas to pending Map tasks based on coefficients
func addReplica(t task, coefficient float32) ([]hash.HashValue, error) {
	mapTaskID := []hash.HashValue{}

	totalSize, err := t.sumSize()
	if err != nil {
		return nil, err
	}
	totalSize = int(float32(totalSize)*coefficient) - totalSize

	for k, _ := range t {
		mapTaskID = append(mapTaskID, k)
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(mapTaskID), func(i, j int) {
		mapTaskID[i], mapTaskID[j] = mapTaskID[j], mapTaskID[i]
	})
	currentSizeSum := 0
	for _, task := range mapTaskID {
		mapTaskID = append(mapTaskID, task)
		currentSizeSum += t[task].FileSize

		// Check if the current sum reaches the target
		if currentSizeSum >= totalSize {
			break
		}
	}
	return mapTaskID, nil
}

//----- M1 of map -----

// assignMapTaskViaSize return the result of assign tasks to receiver via size without consider of sender
func assignMapTaskViaSize(receivers []*Worker, transmitTaskID []hash.HashValue) ([]hash.HashValue, error) {
	//intermediateMapTask WorkerID:taskID assigned tasks to receiver without sender,ideal taskList after map
	intermediateMapTask = make(map[string]hash.HashValue)
	
	averageSize:=len(receivers)/
	//unusedSize WorkerID:totalSize/n
	unusedSize := make(map[string]int)
	//fill up unusedSize
	for _, receiver := range receivers {
		unusedSize[receiver.WorkerID]=
	}

	for _, receiver := range receivers {
		for _, taskID := range transmitTaskID {
			_, receiverGet := receiver.TaskList[taskID]

			if !added && receiverGet {
				t[receiver.Addr] = append(t[receiver.Addr], taskID)
			}
		}
	}
	return t
}

// M1. Bandwidth average assign
// First assign task(assignMapTaskViaSize),then assign sender
func (c *Coordinator) assignMapTaskM1(nonemptyIDList []string, onlineList []*Worker) (mapTaskSet, error) {
	m := make(mapTaskSet)
	mapTasks, err := addReplica(c.allTask, c.ReplicaCoefficient)
	if err != nil {
		return nil, err
	}

	// for _, onlineWorkerID := range onlineList {
	// 	for _, taskID := range mapTasks {
	// 		_, workerGet := onlineWorkerID.TaskList[taskID]

	// 	}
	// }
	return nil
}

// M2. Bandwidth average assign
func (c *Coordinator) assignMapTaskM2(nonemptyIDList []string) mapTaskSet {
	return nil
}

// ----- Reduce -----
// M1.if worker get this task assign(Round-rand assignment)
func (c *Coordinator) assignReduceTaskM1(onlineList []*Worker) reduceTaskSet {
	r := make(reduceTaskSet)
	for _, worker := range onlineList {
		for taskID, _ := range c.allTask {
			_, workerOk := worker.TaskList[taskID]
			_, added := r[string(taskID)]
			if !added && workerOk {
				r[worker.WorkerID] = append(r[worker.WorkerID], taskID)
			}
		}
	}
	return r
}

// M2.Equally distributed according to the number of workers(Average assignment)
func (c *Coordinator) assignReduceTaskM2(onlineList []*Worker) reduceTaskSet {
	return nil
}

// M3.Assign based on worker connections on the basis of 2(RTT assignment )
func (c *Coordinator) assignReduceTaskM3(onlineList []*Worker) reduceTaskSet {
	return nil
}
