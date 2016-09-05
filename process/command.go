package process

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"sync"

	"rde-tech.vir888.com/dev/secretary/secretary.git/parser"
)

const (
	// MULTI mean can run the same command more the once
	MULTI = -1
	// FOREVER mean the command will run one by one and never stop
	FOREVER = 0
)

// Command structure schedule line
type Command struct {
	*sync.RWMutex

	ID          string
	Try         func(parser.SpecSchedule) bool
	cmd         string
	job         func() ([]byte, error)
	timeSetting string
	canRunMulti bool
	mustRemove  bool
	repeat      int
	running     bool
}

func newCmd(id, command, timeSet string, repeat int) (*Command, error) {

	s, err := parser.Parse(timeSet)
	if err != nil {
		return nil, err
	}

	return &Command{
		RWMutex:     &sync.RWMutex{},
		ID:          id,
		canRunMulti: repeat == MULTI,
		mustRemove:  repeat > FOREVER,
		repeat:      repeat,
		cmd:         command,
		job:         createJob(command),
		timeSetting: timeSet,
		Try: func(now parser.SpecSchedule) bool {
			return (now.Second&s.Second) > 0 &&
				(now.Minute&s.Minute) > 0 &&
				(now.Hour&s.Hour) > 0 &&
				(now.Dom&s.Dom) > 0 &&
				(now.Month&s.Month) > 0 &&
				(now.Dow&s.Dow) > 0
		},
		running: false,
	}, nil
}

// CanRunMulti return bool that if command can run over once
func (c *Command) CanRunMulti() bool {
	return c.canRunMulti
}

// MustRemove return if command must count for unregister
func (c *Command) MustRemove() bool {
	return c.mustRemove
}

// IsRunning return command run status
func (c *Command) IsRunning() bool {
	c.RLock()
	defer c.RUnlock()

	return c.running
}

// SetRunning set command running
func (c *Command) SetRunning(fn func() bool) {
	c.Lock()
	defer c.Unlock()
	c.running = fn()
}

// Repeat return Command repeat
func (c *Command) Repeat() int {
	c.RLock()
	defer c.RUnlock()

	return c.repeat
}

func (c *Command) exec() (execOut []byte, err error) {
	if c.MustRemove() && c.Repeat() == 0 {
		return
	}
	c.SetRunning(func() bool {
		return true
	})
	execOut, err = c.job()
	c.SetRunning(func() bool {
		if c.MustRemove() && err == nil {
			c.repeat--
		}
		return false
	})

	return
}

// String convert command to raw string
func (c *Command) String() string {
	return fmt.Sprintf("%s|%d|%s", c.timeSetting, c.repeat, c.cmd)
}

// Hash return sha1(command)
func Hash(cmd string) string {
	hash := sha1.New()
	hash.Write([]byte(cmd))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

var client = &http.Client{}

func createJob(cmd string) func() ([]byte, error) {

	if strings.HasPrefix(cmd, "http://") || strings.HasPrefix(cmd, "https://") {
		return func() ([]byte, error) {
			res, err := client.Get(cmd)
			if err != nil {
				return []byte(""), err
			}
			defer res.Body.Close()

			return ioutil.ReadAll(res.Body)
		}
	}

	return func() ([]byte, error) {
		c := exec.Command("sh", "-c", cmd)
		return c.Output()
	}

}
