package parsnip

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Sender struct {
	registry *Registry
	t        Transport
}

func NewSender(r *Registry, t Transport) *Sender {
	return &Sender{
		registry: r,
		t:        t,
	}
}

func (s *Sender) RunAt(j Job, t time.Time) error {
	name := jobName(j)
	if !s.registry.IsRegisted(j) {
		return errors.Errorf("unregisted job: %s", name)
	}
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(&j)
	if err != nil {
		return errors.Wrap(err, "encode job")
	}
	task := &Task{
		ID:         uuid.New().String(),
		Name:       name,
		Queue:      s.registry.jobs[name].queue,
		Payload:    b.Bytes(),
		StartTime:  t,
		MaxRetries: s.registry.jobs[name].retries,
	}
	return s.t.Enqueue(task)
}

func (s *Sender) RunAfter(j Job, d time.Duration) error {
	return s.RunAt(j, time.Now().Add(d))
}

func (s *Sender) Run(j Job) error {
	return s.RunAt(j, time.Now())
}

// defaults for easier calling semantics

var defaultSender = NewSender(defaultRegistry, &mockTransport{})

func RunAt(j Job, t time.Time) error {
	return defaultSender.RunAt(j, time.Now())
}

func RunAfter(j Job, d time.Duration) error {
	return RunAt(j, time.Now().Add(d))
}

func Run(j Job) error {
	return RunAt(j, time.Now())
}
