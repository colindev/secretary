package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

type Command struct {
	Id      string
	Repeat  int
	Raw     string
	Try     *regexp.Regexp
	Running bool
}

type Process struct {
	Commands map[string]*Command
}

func (p *Process) Recieve(v url.Values) (err error) {

	var repeat int
	var cmd string
	var t *regexp.Regexp

	repeat, err = strconv.Atoi(v.Get("repeat"))

	if err != nil {
		return
	}

	cmd = v.Get("command")
	// TODO: 檢查命令

	t, err = regexp.Compile(v.Get("datetime"))

	if err != nil {
		return
	}

	hash := sha1.New()
	hash.Write([]byte(cmd))
	id := fmt.Sprintf("%x", hash.Sum(nil))

	if _, double := p.Commands[id]; double {
		err = errors.New("重複命令")
		return
	}

	p.Commands[id] = &Command{
		Id:      id,
		Repeat:  repeat,
		Raw:     cmd,
		Try:     t,
		Running: false,
	}

	return
}

func NewProcess() (p *Process) {
	p = &Process{}
	p.Commands = make(map[string]*Command)
	return
}
