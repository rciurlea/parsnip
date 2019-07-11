package parsnip

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"reflect"
)

type Job interface {
	Run(ctx context.Context) (retry bool, err error)
}

// Registry stores "immutable" (after initialization) job information:
// - all queue names
// - queue weights, used by workers; when a worker monitors several queues,
//   it should pick up a job to work on probabilistically using the weights
// - which jobs map to which queues (one job can only be sent on a single queue)
// - job options (retry and backoff policy)
// - job resource requests (named resources each job needs, workers will need to
//   use middleware that provides this, if they want to listen on that queue)

type Registry struct {
	jobs map[string]bool
}

func NewRegistry() *Registry {
	return &Registry{
		jobs: map[string]bool{},
	}
}

func (r *Registry) Add(j Job) {
	gob.Register(j)
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(j); err != nil {
		panic(fmt.Sprintf("job not serializable: %v", err))
	}
	typ := reflect.TypeOf(j)
	jobName := typ.PkgPath() + "." + typ.Name()
	if r.jobs[jobName] {
		panic(fmt.Sprintf("job already in registry: %s", jobName))
	}
	r.jobs[jobName] = true
}

func (r *Registry) Print() {
	for j := range r.jobs {
		fmt.Println(j)
	}
}
