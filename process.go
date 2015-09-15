package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Command struct {
	Id      string
	Repeat  int
	Raw     string
	Try     func(string) bool
	Running bool
}

type Process struct {
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

	re, err = parseDatetime(v.Get("datetime"))

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
		Try:     func(now string) bool { m := re.FindString(now); return len(m) > 0 },
		Running: false,
	}

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

	any := "[0-9]{1,2}"
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
	re := fmt.Sprintf("[0-9]{4}-%s-%sT%s:%s:%s+00:00", m, d, h, i, s)
	fmt.Println("正規式:", re)

	r, e = regexp.Compile(re)

	return
}
