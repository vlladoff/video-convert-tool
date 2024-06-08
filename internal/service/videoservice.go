package service

import (
	"context"
	"github.com/vlladoff/video-convert-tool/internal/config"
	"github.com/vlladoff/video-convert-tool/internal/entity"
	"github.com/vlladoff/video-convert-tool/internal/repository"
	"github.com/vlladoff/video-convert-tool/internal/storage"
	"github.com/vlladoff/video-convert-tool/internal/workerpool"
)

type VideoService struct {
	taskRepo  *repository.TaskRepository
	videoRepo *repository.VideoRepository
	wp        *workerpool.WorkerPool
}

func NewVideoService(cfg *config.Config, wp *workerpool.WorkerPool, storage storage.Storage) (*VideoService, func()) {
	taskRepo, closeFunc := repository.NewTaskRepository(cfg, wp.Workers)
	videoRepo := repository.NewVideoRepository(storage)
	return &VideoService{taskRepo: taskRepo, videoRepo: videoRepo, wp: wp}, closeFunc
}

func (vs *VideoService) StartConsumingTasks(ctx context.Context) {
	go vs.taskRepo.StartReading(ctx)
	go vs.DoneTasks(ctx)

	for t := range vs.taskRepo.GetTasksChan() {
		vs.wp.AddTask(&t)
	}

	vs.wp.Wait()
}

func (vs *VideoService) DoneTasks(ctx context.Context) {
	go vs.taskRepo.StartWriting(ctx)

	for done := range vs.wp.DoneChan {
		taskDone := entity.ConvertVideoTaskDone{
			ID:     done.Id,
			Status: done.Status,
		}

		err := vs.taskRepo.WriteCompletedTask(ctx, taskDone)
		// todo log
		if err != nil {
			return
		}
	}
}
