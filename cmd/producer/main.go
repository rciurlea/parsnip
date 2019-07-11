package main

import (
	"github.com/rciurlea/parsnip"
	"github.com/rciurlea/parsnip/cmd/jobs"
	j2 "github.com/rciurlea/parsnip/cmd/otherjobs/jobs"
)

func main() {
	r := parsnip.NewRegistry()
	r.Add(jobs.A{})
	r.Add(j2.A{})
	r.Print()
}
