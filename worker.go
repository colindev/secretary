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
	defer w.ticker.Stop()

	for t := range ct {
		time_string := t.Format(time.RFC3339)
		p.Each(func(cmd *Command, id string) error {
			if cmd.Try(time_string) {
				go func(script string, t string) {
					c := exec.Command("sh", "-c", script)
					out, err := c.Output()
					fmt.Printf("\033[33m[ %s ] exec: %s, out: %+v, err: %+v\033[m\n", t, script, out, err)
				}(cmd.Cmd, time_string)
			}

			return nil
		})
	}
}

func NewWork(i int64) *Worker {
	return &Worker{interval: i}
}
