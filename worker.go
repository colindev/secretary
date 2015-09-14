package main

import (
	"fmt"
	"time"
)

type Worker struct {
	interval int64
}

func (w Worker) Run(i int64) {

	w.interval = i

	ct := time.NewTicker(time.Second * time.Duration(i)).C

	go func() {
		for t := range ct {
			fmt.Println("time.Ticker.C: ", processes.Commands, t)
		}
	}()
}

var worker = Worker(Worker{})
