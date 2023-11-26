package mr

import "main/utils"

type Task map[utils.HashValue]FileMeta

type TransmitTaskSet map[string][]utils.HashValue //WorkerAddr:[]TaskID

type SingleTransmitTask struct {
	TaskID utils.HashValue
	FMeta  FileMeta
	FData  []byte
}

func MakeTransmitTaskSet(receivers []*Worker, transmitTaskID []utils.HashValue) TransmitTaskSet {
	t := make(TransmitTaskSet)
	for _, receiver := range receivers {
		for _, taskID := range transmitTaskID {
			if _, ok := receiver.TaskList[taskID]; ok == true {
				t[receiver.Addr] = append(t[receiver.Addr], taskID)
			}
		}
	}
	return t
}
