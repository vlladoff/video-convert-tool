package worker_pool

import (
	"sync"
)

type Task interface {
	Process(workerId int)
}

type TaskMeta struct {
	Id int
}

type WorkerPool struct {
	Workers   int
	TasksChan chan Task
	Wg        sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		Workers:   workers,
		TasksChan: make(chan Task, workers),
	}
}

func (wp *WorkerPool) StartWorkers() {
	for i := 0; i < wp.Workers; i++ {
		wp.Wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.TasksChan <- task
}

func (wp *WorkerPool) worker(id int) {
	defer wp.Wg.Done()

	for task := range wp.TasksChan {
		task.Process(id)
	}
}
