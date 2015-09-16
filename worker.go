package main

import (
	"fmt"
	"os/exec"
	"time"
)

type Worker struct {
	interval int64
	ticker   *time.Ticker
}

func (w *Worker) Run(p *Process) {

	w.ticker = time.NewTicker(time.Second * time.Duration(w.interval))
	ct := w.ticker.C

	go func() {
		for t := range ct {
			time_string := t.Format(time.RFC3339)
			fmt.Println("time.Ticker.C to RFC3339:", time_string)
			p.Each(func(cmd *Command, id string) error {
				if cmd.Try(time_string) {
					go func(script string) {
						c := exec.Command("sh", "-c", script)
						out, err := c.Output()
						fmt.Printf("exec: %s, out: %+v, err: %+v\n\n", script, out, err)
					}(cmd.Cmd)
				}

				return nil
			})
		}
	}()
}

func (w *Worker) Stop() {
	w.ticker.Stop()
}

func NewWork(i int64) *Worker {
	return &Worker{interval: i}
}
