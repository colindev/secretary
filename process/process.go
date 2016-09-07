package process

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"rde-tech.vir888.com/dev/secretary/secretary.git/parser"
)

var (
	// ErrCommandDuplicate mean has receive the same command
	ErrCommandDuplicate = errors.New("[process] Command duplicate")
)

type running struct {
	*sync.RWMutex
	list map[string]bool
}

func (r *running) add(str string) {
	r.Lock()
	defer r.Unlock()
	r.list[str] = true
}

func (r *running) del(str string) {
	r.Lock()
	defer r.Unlock()
	delete(r.list, str)
}

func (r *running) all() []string {
	r.RLock()
	defer r.RUnlock()
	ret := []string{}
	for str := range r.list {
		ret = append(ret, str)
	}
	return ret
}

// Process cache the running command
type Process struct {
	*sync.RWMutex
	*sync.WaitGroup
	*log.Logger
	schedules map[string]*Command
	r         *running
	tk        *time.Ticker
}

// New return process instance
func New(log *log.Logger) *Process {
	return &Process{
		RWMutex:   &sync.RWMutex{},
		WaitGroup: &sync.WaitGroup{},
		Logger:    log,
		schedules: map[string]*Command{},
		r: &running{
			RWMutex: &sync.RWMutex{},
			list:    map[string]bool{},
		},
	}
}

// Receive create command and cache it
func (p *Process) Receive(command, timeSet string, repeat int) (id string, err error) {

	id = Hash(command)

	p.Lock()
	defer p.Unlock()
	if _, double := p.schedules[id]; double {
		err = ErrCommandDuplicate
		return
	}

	p.schedules[id], err = newCmd(id, command, timeSet, repeat)
	return
}

// Revoke delete exists command from process
func (p *Process) Revoke(id string) {
	p.Lock()
	delete(p.schedules, id)
	p.Unlock()
}

// Each iterate all exists commands and invoke callback
func (p *Process) Each(f func(*Command) error) {
	p.RLock()
	for _, c := range p.schedules {
		if err := f(c); err != nil {
			break
		}
	}
	p.RUnlock()
}

// Backup dump all exists commands to specific file
func (p *Process) Backup(file string, s time.Duration) {

	p.tk = time.NewTicker(s)

	go func() {
		for _ = range p.tk.C {
			p.Printf("[backup] %v\n", p.backupTo(file))
		}
	}()
}

func (p *Process) backupTo(file string) (err error) {
	// TODO 研究有沒有更順的開檔方式
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.WriteString(strings.Join(p.dump(func(c *Command) string { return "" }), "\n"))
	if err != nil {
		return
	}

	return f.Sync()
}

// Find commands match schedule
func (p *Process) Find(s parser.SpecSchedule) []*Command {

	ret := []*Command{}

	p.Each(func(c *Command) error {
		if c.Try(s) {
			ret = append(ret, c)
		}

		return nil
	})

	return ret
}

// Exec commands
func (p *Process) Exec(t time.Time, cs []*Command) {

	p.Lock()
	defer p.Unlock()

	for _, c := range cs {
		if c.IsRunning() && !c.CanRunMulti() {
			continue
		}
		c.SetRunning(func() bool { return true })
		// last time
		if c.Repeat() == 1 {
			delete(p.schedules, c.ID)
		}

		go func(cmd *Command) {
			p.Add(1)
			ind := fmt.Sprintf("%s: %s", t, cmd.String())
			p.r.add(ind)
			out, err := cmd.exec()
			p.r.del(ind)
			p.Printf("\033[33m %s => %+v\n%s\033[m\n", cmd.String(), err, out)
			p.Done()
		}(c)
	}
}

// Running return running commands with time
func (p *Process) Running() []string {
	return p.r.all()
}

// Stop backup ticker
func (p *Process) Stop() {
	if p.tk != nil {
		p.tk.Stop()
	}
}

func (p *Process) dump(f func(c *Command) string) []string {
	s := []string{}

	p.Each(func(c *Command) (err error) {

		s = append(s, f(c)+c.String())

		return
	})

	return s
}
