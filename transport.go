package parsnip

import (
	"time"
)

type Task struct {
	ID         string
	Name       string
	Queue      string
	Payload    []byte
	StartTime  time.Time
	Retries    uint
	MaxRetries uint
}

type Transport interface {
	EnsureSchema(r *Registry) error
	Enqueue(t *Task) error
	Dequeue(workerID string) (*Task, error)
	Done(taskID string) error
}

type mockTransport struct {
	tasks []*Task
}

func (mt *mockTransport) EnsureSchema(r *Registry) error {
	return nil
}

func (mt *mockTransport) Enqueue(t *Task) error {
	mt.tasks = append(mt.tasks, t)
	return nil
}

func (mt *mockTransport) Dequeue(workerID string) (*Task, error) {
	return nil, nil
}

func (mt *mockTransport) Done(taskID string) error {
	return nil
}
