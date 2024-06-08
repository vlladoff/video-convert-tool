package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/vlladoff/video-convert-tool/internal/config"
	"github.com/vlladoff/video-convert-tool/internal/entity"
)

const (
	topicName     = "vct_task"
	topicDoneName = "vct_task_done"
	maxBytes      = 10e6 // 10 mb
)

type TaskRepository struct {
	reader *kafka.Reader
	writer *kafka.Writer
	tasks  chan entity.ConvertVideoTask
	done   chan entity.ConvertVideoTaskDone
}

func NewTaskRepository(cfg *config.Config, maxTasks int) (*TaskRepository, func()) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaAddress},
		GroupID:  cfg.GroupId,
		Topic:    topicName,
		MaxBytes: maxBytes,
	})

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaAddress),
		Topic:    topicDoneName,
		Balancer: &kafka.LeastBytes{},
	}

	tasks := make(chan entity.ConvertVideoTask, maxTasks)
	doneTasks := make(chan entity.ConvertVideoTaskDone, maxTasks)

	return &TaskRepository{reader: reader, writer: writer, tasks: tasks, done: doneTasks}, func() {
		err := reader.Close()
		// todo log
		if err != nil {
			return
		}
		err = writer.Close()
		if err != nil {
			return
		}
		close(tasks)
		close(doneTasks)
	}
}

func (r *TaskRepository) StartReading(ctx context.Context) {
	for {
		m, err := r.reader.ReadMessage(ctx)
		// todo log
		if err != nil {
			break
		}

		var newTask entity.ConvertVideoTask
		err = json.Unmarshal(m.Value, &newTask)
		if err != nil {
			return
		}

		r.tasks <- newTask

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}

func (r *TaskRepository) StartWriting(ctx context.Context) {
	for done := range r.done {
		err := r.WriteCompletedTask(ctx, done)
		// todo log
		if err != nil {
			return
		}
	}
}

func (r *TaskRepository) WriteCompletedTask(ctx context.Context, t entity.ConvertVideoTaskDone) error {
	taskJSON, err := json.Marshal(t)
	// todo log
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %s\n", err)
		return err
	}

	err = r.writer.WriteMessages(ctx, kafka.Message{
		Value: taskJSON,
	})

	return err
}

func (r *TaskRepository) GetTasksChan() <-chan entity.ConvertVideoTask {
	return r.tasks
}
