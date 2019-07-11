package parsnip

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/pkg/errors"
)

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

// default stuff, for quicker use (similar to net/http)
var defaultRegistry = NewRegistry()

func DefaultRegistry() *Registry {
	return defaultRegistry
}

func Register(j Job, opts ...JobOption) error {
	return defaultRegistry.Register(j, opts...)
}

func MustRegister(j Job, opts ...JobOption) {
	defaultRegistry.MustRegister(j, opts...)
}
