package main

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Process struct {
	Flag     string `json:"flag"`
	Command  string `json:"command"`
	Datetime string `json:"datetime"`
}

var processes = make(map[string]map[string]string)

func ProcessesParse(v url.Values) (*Process, error) {

	var err error
	p := &Process{
		Flag:     v.Get("flag"),
		Command:  v.Get("command"),
		Datetime: v.Get("datetime"),
	}

	switch p.Flag {
	case "once":
	case "repeat":
	case "unique":
		break
	default:
		p = nil
		err = errors.New("不允許的旗標")
	}

	return p, err

}

func ProcessesReceiver(i int64) func(*Process) error {

	ct := time.NewTicker(time.Second * time.Duration(i)).C

	go func() {
		for t := range ct {
			fmt.Println(t, processes)
		}
	}()

	return func(p *Process) error {

		k := p.Flag
		command := p.Command
		datetime := p.Datetime

		// TODO: check command

		// TODO: check datetime

		cmds, has_map := processes[k]

		// auto create map
		if !has_map {
			processes[k] = make(map[string]string)
		}

		if _, duble := cmds[command]; duble {
			return errors.New("命令重複")
		}

		processes[k][command] = datetime

		return nil
	}

}
