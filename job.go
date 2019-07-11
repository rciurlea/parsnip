package parsnip

import (
	"context"
	"reflect"
)

// Job must be implemented by types that are to be queued for execution
type Job interface {
	Run(ctx context.Context) (retry bool, err error)
}

func jobName(j Job) string {
	typ := reflect.TypeOf(j)
	return typ.PkgPath() + "." + typ.Name()
}

type jobParams struct {
	queue   string
	deps    []string
	retries uint
}

type JobOption func(*jobParams)

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
