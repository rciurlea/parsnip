package parsnip

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/pkg/errors"
)

// Registry stores "immutable" (after initialization) job information:
// - all queue names
// - queue weights, used by workers; when a worker monitors several queues,
//   it should pick up a job to work on probabilistically using the weights
// - which jobs map to which queues (one job can only be sent on a single queue)
// - job options (retry and backoff policy)
// - job resource requests (named resources each job needs, workers will need to
//   use middleware that provides this, if they want to listen on that queue)

const defaultQueue = "default"

type Registry struct {
	queueDeps map[string][]string
	jobs      map[string]*jobParams
}

func NewRegistry() *Registry {
	return &Registry{
		queueDeps: map[string][]string{
			defaultQueue: []string{},
		},
		jobs: map[string]*jobParams{},
	}
}

func (r *Registry) Print() {
	for j, opts := range r.jobs {
		fmt.Printf("%s: %+v\n", j, opts)
	}
}

func (r *Registry) Register(j Job, opts ...JobOption) error {
	params := &jobParams{
		queue:   defaultQueue,
		retries: 1,
		deps:    []string{},
	}
	for _, o := range opts {
		o(params)
	}
	gob.Register(j)
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(j); err != nil {
		return errors.Wrap(err, "job not serializable")
	}
	if _, ok := r.jobs[jobName(j)]; ok {
		return errors.Errorf("%s: job already registered", jobName)
	}
	r.jobs[jobName(j)] = params
	return nil
}

func (r *Registry) MustRegister(j Job, opts ...JobOption) {
	if err := r.Register(j, opts...); err != nil {
		panic(err)
	}
}
