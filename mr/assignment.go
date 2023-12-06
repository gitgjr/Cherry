package mr

// M1. Average assign without bandwidth
func (c *Coordinator) assignMapTaskM1(nonemptyIDList []string) mapTaskSet {
	return nil
}

// M2. Bandwidth average assign
func (c *Coordinator) assignMapTaskM2(nonemptyIDList []string) mapTaskSet {
	return nil
}

// M1.if worker get this task assign(Random assignment)
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
