package main

import (
	"net/url"
	"regexp"
	"strconv"
)

type Process struct {
	Repeat  int
	Command string
	Try     *regexp.Regexp
	Running bool
}

type Processes struct {
	Commands []*Process
	interval int64
}

func (x Processes) Recieve(v url.Values) (err error) {

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

	p := &Process{Repeat: repeat, Command: cmd, Try: t, Running: false}
	x.Commands = append(x.Commands, p)

	return
}

var processes = Processes(Processes{})
