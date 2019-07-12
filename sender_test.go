package parsnip

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSender(t *testing.T) {
	j := MockJob{}
	r := NewRegistry()
	r.Register(j)
	mt := &mockTransport{}
	s := NewSender(r, mt)
	start := time.Date(2019, 7, 12, 23, 28, 0, 0, time.UTC)
	err := s.RunAt(j, start)
	assert.NoError(t, err)
	assert.Len(t, mt.tasks, 1)
	task := mt.tasks[0]
	assert.Equal(t, jobName(j), task.Name)
	assert.Equal(t, r.jobs[jobName(j)].queue, task.Queue)
	assert.Equal(t, r.jobs[jobName(j)].retries, task.MaxRetries)
	assert.Equal(t, start, task.StartTime)
}
