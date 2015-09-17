package main

import (
	"flag"
)

func main() {

	addr := flag.String("addr", "", "開發測試用")
	interval := flag.Int64("interval", 10, "排程掃描間隔")
	f_conf := flag.String("schedule", "", "排程設定檔")
	f_backup := flag.String("backup", "/tmp/schedule.backup", "排程備份位置")

	flag.Parse()

	worker := NewWork(*interval)
	process := NewProcess(*f_conf)

	// TODO: 評估要不要拿掉
	if *addr != "" {
		ListenAndServe(*addr, process)
	}

	// 背景備份
	process.Backup(*f_backup, 60)
	worker.Run(process)

}
