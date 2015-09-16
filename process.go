package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Command struct {
	Id      string
	Repeat  int
	Cmd     string
	Try     func(string) bool
	Running bool
	Raw     string
}

type Process struct {
	// TODO: 壞味道,Worker 必須自己迭代 Commands
	Commands map[string]*Command
}

func (p *Process) Recieve(v url.Values) (err error) {

	var repeat int
	var cmd string
	var re *regexp.Regexp

	repeat, err = strconv.Atoi(v.Get("repeat"))

	if err != nil {
		return
	}

	cmd = v.Get("command")
	// TODO: 檢查命令

	dt := v.Get("datetime")
	re, err = parseDatetime(dt)

	if err != nil {
		return
	}

	hash := sha1.New()
	hash.Write([]byte(cmd))
	id := fmt.Sprintf("%x-%s", hash.Sum(nil), cmd)

	if _, double := p.Commands[id]; double {
		err = errors.New("重複命令")
		return
	}

	// TODO: 設計個比較優雅的方式傳出
	fmt.Println("\033[32mregexp:", re.String(), "\033[m")

	p.Commands[id] = &Command{
		Id:      id,
		Repeat:  repeat,
		Cmd:     cmd,
		Try:     func(now string) bool { m := re.FindString(now); return len(m) > 0 },
		Running: false,
		Raw:     id + "|" + dt + "|" + strconv.Itoa(repeat) + "|" + cmd,
	}

	return
}

func (p *Process) Revoke(id string) {
	delete(p.Commands, id)
}

func (p *Process) Each(f func(*Command, string) error) {
	for id, c := range p.Commands {
		if err := f(c, id); err != nil {
			break
		}
	}
}

func (p *Process) Backup(file string) (err error) {
	f, err := os.Create(file)
	defer f.Close()
	// TODO: 開檔失敗無法備份
	if err != nil {
		return
	}

	p.Each(func(c *Command, id string) (err error) {
		_, err = f.WriteString(c.Raw + "\n")

		return
	})

	f.Sync()

	return
}

func NewProcess() (p *Process) {
	p = &Process{}
	p.Commands = make(map[string]*Command)
	return
}

func parseDatetime(s string) (r *regexp.Regexp, e error) {

	// 秒 分 時 日 月 星期
	arr := strings.Split(s, " ")
	if 6 != len(arr) {
		e = errors.New("格式為: 秒 分 時 日 月 星期")
		return
	}

	any := "[0-9]{2}"
	for n := 0; n <= 5; n++ {

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
