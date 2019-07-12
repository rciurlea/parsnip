package parsnip

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockJob struct {
}

func (mj MockJob) Run(ctx context.Context) (retry bool, err error) {
	return false, nil
}

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	j := MockJob{}
	assert.False(t, r.IsRegisted(j))
	jName := jobName(j)
	err := r.Register(j, OptQueue("qq"), OptRetries(10), OptDeps("1", "2"))
	assert.NoError(t, err)
	err = r.Register(j)
	assert.Error(t, err) // error on repeated registration
	assert.True(t, r.IsRegisted(j))
	assert.Contains(t, r.jobs, jName)
	assert.Equal(t, "qq", r.jobs[jName].queue)
	assert.Equal(t, uint(10), r.jobs[jName].retries)
	assert.Contains(t, r.jobs[jName].deps, "1")
	assert.Contains(t, r.jobs[jName].deps, "2")
	assert.Contains(t, r.queueDeps, "qq")
	assert.Contains(t, r.queueDeps["qq"], "1")
	assert.Contains(t, r.queueDeps["qq"], "2")
}
