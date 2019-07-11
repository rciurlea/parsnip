package parsnip

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
)

type Sender struct {
	registry *Registry
	buf      *bytes.Buffer
	enc      *gob.Encoder
}

func NewSender(r *Registry) *Sender {
	b := &bytes.Buffer{}
	return &Sender{
		registry: r,
		buf:      b,
		enc:      gob.NewEncoder(b),
	}
}

func (s *Sender) Run(j Job) error {
	return s.enc.Encode(&j)
}

func (s *Sender) Decode() error {
	var j Job
	err := gob.NewDecoder(s.buf).Decode(&j)
	if err != nil {
		return err
	}
	j.Run(context.Background())
	log.Println(jobName(j))
	return nil
}

// defaults for easier calling semantics

var defaultSender = NewSender(defaultRegistry)

func Run(j Job) error {
	return defaultSender.Run(j)
}

func Decode() error {
	return defaultSender.Decode()
}
