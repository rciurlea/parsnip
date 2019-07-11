package main

import (
	"github.com/rciurlea/parsnip"
	"github.com/rciurlea/parsnip/cmd/jobs"
)

func main() {
	r := parsnip.NewRegistry()
	r.MustRegister(jobs.A{},
		parsnip.OptQueue("meh"),
		parsnip.OptRetries(5),
		parsnip.OptDeps("db", "log"),
	)
	r.Print()
}
