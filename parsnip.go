package parsnip

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// Job must be implemented by types that are to be queued for execution
type Job interface {
	Run(ctx context.Context) (retry bool, err error)
}

type jobParams struct {
	queue   string
	deps    []string
	retries uint
}

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

type JobOption func(*jobParams)

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
	typ := reflect.TypeOf(j)
	jobName := typ.PkgPath() + "." + typ.Name()
	if _, ok := r.jobs[jobName]; ok {
		return errors.Errorf("%s: job already registered", jobName)
	}
	r.jobs[jobName] = params
	return nil
}

func (r *Registry) MustRegister(j Job, opts ...JobOption) {
	if err := r.Register(j, opts...); err != nil {
		panic(err)
	}
}

func OptQueue(q string) JobOption {
	return func(p *jobParams) {
		p.queue = q
	}
}

func OptRetries(r uint) JobOption {
	return func(p *jobParams) {
		p.retries = r
	}
}

func OptDeps(deps ...string) JobOption {
	return func(p *jobParams) {
		p.deps = deps
	}
}
