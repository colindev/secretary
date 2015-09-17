package main

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Command struct {
	Id          string
	Repeat      int
	Cmd         string
	TimeSetting string
	Try         func(string) bool
	Running     bool
}

func (c *Command) Raw() string {

	i := strconv.Itoa(c.Repeat)

	return c.TimeSetting + "|" + i + "|" + c.Cmd
}

type Process struct {
	Schedules map[string]*Command
}

func (p *Process) Recieve(repeat int, command string, time_set string) (err error) {

	// TODO: 評估檢查 command 命令字串 或是 執行失敗移除

	re, err := parseDatetime(time_set)

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
	fmt.Println("\033[32mregexp:", re.String(), "\033[m")

	p.Schedules[id] = &Command{
		Id:          id,
		Repeat:      repeat,
		Cmd:         command,
		TimeSetting: time_set,
		Try:         func(now string) bool { m := re.FindString(now); return len(m) > 0 },
		Running:     false,
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

			arr := strings.Split(text, "|")

			if len(arr) != 3 {
				continue
			}

			repeat, _ := strconv.Atoi(arr[1])
			command := arr[2]
			time_set := arr[0]
			p.Recieve(repeat, command, time_set)

			fmt.Println("schedule:", text)
		}
	}

	return
}

func parseDatetime(s string) (r *regexp.Regexp, e error) {

	// 秒 分 時 日 月
	arr := strings.Split(s, " ")
	if 5 != len(arr) {
		e = errors.New("格式為: 秒 分 時 日 月")
		return
	}

	any := "[0-9]{2}"
	for n := 0; n < len(arr); n++ {

		sub := strings.Split(arr[n], ",")

		for m := 0; m < len(sub); m++ {
			if "*" == sub[m] {
				sub[m] = any
				continue
			}
			num, err := strconv.Atoi(sub[m])
			if err != nil {
				e = errors.New("只能使用*或數字 如 1 或是 * 或是 1,6")
				return
			}
			sub[m] = fmt.Sprintf("%02d", num)
		}

		arr[n] = "(?:" + strings.Join(sub, "|") + ")"
	}

	m, d, h, i, s := arr[4], arr[3], arr[2], arr[1], arr[0]

	// use RFC3339
	re := fmt.Sprintf("[0-9]{4}-%s-%sT%s:%s:%s\\+%s:%s", m, d, h, i, s, any, any)

	r, e = regexp.Compile(re)

	return
}
