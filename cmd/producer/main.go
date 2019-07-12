package main

import (
	"log"

	"github.com/rciurlea/parsnip"
	"github.com/rciurlea/parsnip/cmd/jobs"
)

func main() {
	parsnip.MustRegister(jobs.A{},
		parsnip.OptQueue("meh"),
		parsnip.OptRetries(5),
		parsnip.OptDeps("db", "log"),
	)
	err := parsnip.Run(jobs.A{Something: "heh!"})
	if err != nil {
		log.Fatal(err)
	}
	err = parsnip.Run(jobs.A{Something: "cool!"})
	if err != nil {
		log.Fatal(err)
	}
}
