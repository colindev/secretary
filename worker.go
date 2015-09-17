package main

import (
	"fmt"
	"os/exec"
	"time"
)

type Worker struct {
	interval int64
	ticker   *time.Ticker
}

func (w *Worker) Run(p *Process) {

	w.ticker = time.NewTicker(time.Second * time.Duration(w.interval))
	ct := w.ticker.C
	defer w.ticker.Stop()

	for t := range ct {
		time_string := t.Format(time.RFC3339)
		p.Each(func(cmd *Command, id string) (err error) {

			// TODO: 壞味道,魔術數字
			// 跳過僅允許同一時間只能跑一條的程序
			// 換句話說 只有設定 -1 才有機會在同一直間跑兩條以上相同程序
			if cmd.Running && cmd.Repeat != -1 {
				return
			}

			if cmd.Try(time_string) {
				go func(cmd *Command, t string) {

					// 砍掉計數型且已經歸零的程序
					if need_revoke := cmd.Repeat > 0; need_revoke {
						cmd.Repeat--
						if 0 == cmd.Repeat {
							p.Revoke(cmd.Id)
						}
					}

					cmd.Running = true
					c := exec.Command("sh", "-c", cmd.Cmd)
					out, err := c.Output()
					cmd.Running = false

					fmt.Printf("\033[33m[ %s ] exec: %s, out: %+v, err: %+v\033[m\n", t, cmd.Cmd, out, err)
					fmt.Println(cmd.Raw())

				}(cmd, time_string)
			}

			return
		})
	}
}

func NewWork(i int64) *Worker {
	return &Worker{interval: i}
}
