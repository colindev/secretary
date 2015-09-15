package main

import (
	"fmt"
	"time"
)

type Worker struct {
	interval int64
}

func (w *Worker) Run(p *Process) {

	ct := time.NewTicker(time.Second * time.Duration(w.interval)).C

	go func() {
		for t := range ct {
			time_string := t.Format(time.RFC3339)
			fmt.Println("time.Ticker.C to RFC3339:", time_string)
			for _, cmd := range p.Commands {
				fmt.Printf("%+v\n", cmd)
				fmt.Printf("%+v\n", cmd.Try(time_string))
			}
		}
	}()
}

func NewWork(i int64) *Worker {
	return &Worker{interval: i}
}
