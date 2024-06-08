package workerpool

import (
	"sync"
)

type Task interface {
	Process(workerId int) (bool, int)
}

type DoneChan struct {
	Id     int
	Status bool
}

type WorkerPool struct {
	Workers   int
	TasksChan chan Task
	DoneChan  chan DoneChan
	Wg        sync.WaitGroup
}

func NewWorkerPool(workers int) (*WorkerPool, func()) {
	taskChan := make(chan Task, workers)
	doneChan := make(chan DoneChan, workers)

	return &WorkerPool{
			Workers:   workers,
			TasksChan: taskChan,
			DoneChan:  doneChan,
		}, func() {
			close(taskChan)
			close(doneChan)
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
		status, id := task.Process(id)
		wp.DoneChan <- DoneChan{
			Id:     id,
			Status: status,
		}
	}
}

func (wp *WorkerPool) Wait() {
	wp.Wg.Wait()
}
