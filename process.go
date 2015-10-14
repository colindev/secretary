package main

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Command struct {
	Id          string
	Repeat      int
	Cmd         string
	TimeSetting string
	Try         func(*SpecSchedule) bool
	Running     bool
}

func (c *Command) Raw() string {

	i := strconv.Itoa(c.Repeat)

	return c.TimeSetting + "|" + i + "|" + c.Cmd
}

type Process struct {
	Schedules map[string]*Command
}

func (p *Process) Receive(repeat int, command string, time_set string) (err error) {

	// TODO: 評估檢查 command 命令字串 或是 執行失敗移除

	s, err := Parse(time_set)

	if err != nil {
		return
	}

	hash := sha1.New()
	hash.Write([]byte(command))
	id := fmt.Sprintf("%x", hash.Sum(nil))

	if _, double := p.Schedules[id]; double {
		err = errors.New("重複命令")
		return
	}

	// TODO: 設計個比較優雅的方式傳出
	fmt.Printf("\033[32mregexp: %+v \033[m\n", s)

	p.Schedules[id] = &Command{
		Id:          id,
		Repeat:      repeat,
		Cmd:         command,
		TimeSetting: time_set,
		Try: func(now *SpecSchedule) bool {
			return (now.Second&s.Second) > 0 &&
				(now.Minute&s.Minute) > 0 &&
				(now.Hour&s.Hour) > 0 &&
				(now.Dom&s.Dom) > 0 &&
				(now.Month&s.Month) > 0 &&
				(now.Dow&s.Dow) > 0
		},
		Running: false,
	}

	return
}

func (p *Process) Revoke(id string) {
	delete(p.Schedules, id)
}

func (p *Process) Each(f func(*Command, string) error) {
	for id, c := range p.Schedules {
		if err := f(c, id); err != nil {
			break
		}
	}
}

func (p *Process) dump(f func(c *Command) string) string {
	s := ""

	p.Each(func(c *Command, id string) (err error) {

		s = s + f(c) + c.Raw() + "\n"

		return
	})

	return s
}

func (p *Process) Backup(file string, i int) <-chan error {

	ce := make(chan error)

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(i))
		ct := ticker.C
		defer ticker.Stop()

		for _ = range ct {
			f, err := os.Create(file)
			defer f.Close()
			// TODO: 開檔失敗無法備份
			if err != nil {
				return
			}

			_, err = f.WriteString(p.dump(func(c *Command) string { return "" }))

			f.Sync()
		}
	}()
	return ce
}

func NewProcess(conf string) (p *Process) {
	p = &Process{}
	p.Schedules = make(map[string]*Command)

	if f, err := os.Open(conf); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			text := scanner.Text()
			if "" == text {
				continue
			}

			// 跳過 # 開頭的註解
			if strings.HasPrefix(text, "#") {
				continue
			}

			// 重新切割字串
			arr := strings.Split(text, "|")

			if len(arr) != 3 {
				continue
			}

			repeat, _ := strconv.Atoi(arr[1])
			command := arr[2]
			time_set := arr[0]
			p.Receive(repeat, command, time_set)

			fmt.Println("schedule:", text)
		}
	}

	return
}
