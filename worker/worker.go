package worker

import (
	"time"

	"rde-tech.vir888.com/dev/secretary/secretary.git/parser"
)

// Worker handle commands
type Worker struct {
	interval int64
	ticker   *time.Ticker
}

// New return Worker instance
func New(i int64) *Worker {
	return &Worker{
		interval: i,
	}
}

// Run start a time.Ticker
func (w *Worker) Run(fn func(time.Time, parser.SpecSchedule)) {

	w.ticker = time.NewTicker(time.Second * time.Duration(w.interval))

	for t := range w.ticker.C {
		s := parser.SpecSchedule{
			Second: 1 << uint(t.Second()),
			Minute: 1 << uint(t.Minute()),
			Hour:   1 << uint(t.Hour()),
			Dom:    1 << uint(t.Day()),
			Month:  1 << uint(t.Month()),
			Dow:    1 << uint(t.Weekday()),
		}
		fn(t, s)
	}
}

// Stop worker ticker and wait all running command end
func (w *Worker) Stop() {
	if w.ticker != nil {
		w.ticker.Stop()
	}
}
