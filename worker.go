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
			fmt.Println("time.Ticker.C:", t)
			for _, cmd := range p.Commands {
				fmt.Printf("%+v\n", cmd)
			}
		}
	}()
}

func NewWork(i int64) *Worker {
	return &Worker{interval: i}
}
