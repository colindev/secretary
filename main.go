package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"rde-tech.vir888.com/dev/secretary/secretary.git/parser"
	"rde-tech.vir888.com/dev/secretary/secretary.git/process"
	"rde-tech.vir888.com/dev/secretary/secretary.git/worker"
)

var (
	// Version of build
	Version string
	// CompileDate of build
	CompileDate string
	// Env collect flag args
	Env = struct {
		path           string
		AdminAPIAddr   string
		AdminAPIPrefix string
		Interval       int64
		Schedule       string
		Backup         string
		BackupInterval int
	}{}

	prc *process.Process
)

func init() {

	flag.StringVar(&Env.path, "env", "", "env filename")
	flag.StringVar(&Env.AdminAPIAddr, "api.addr", "", "開發測試用(預計廢除)")
	flag.StringVar(&Env.AdminAPIPrefix, "api.prefix", "", "API 路徑前綴")
	flag.Int64Var(&Env.Interval, "interval", 0, "排程掃描間隔")
	flag.StringVar(&Env.Schedule, "schedule", "", "排程設定檔")
	flag.StringVar(&Env.Backup, "backup", "", "排程備份位置(僅在啟動API時才有用)")
	flag.IntVar(&Env.BackupInterval, "backup.interval", 0, "備份間隔秒數(僅在啟動API時才有用)")

	var showVersion = flag.Bool("v", false, "display current version")
	flag.Parse()

	if *showVersion {
		fmt.Println(Version, CompileDate)
		os.Exit(0)
	}

	if err := godotenv.Load(Env.path); err == nil {
		if Env.AdminAPIAddr == "" {
			Env.AdminAPIAddr = os.Getenv("ADMIN_API_ADDR")
		}
		if Env.AdminAPIPrefix == "" {
			Env.AdminAPIPrefix = os.Getenv("ADMIN_API_PREFIX")
		}
		if Env.Interval == 0 {
			Env.Interval, _ = strconv.ParseInt(os.Getenv("INTERVAL"), 10, 64)
		}
		if Env.Schedule == "" {
			Env.Schedule = os.Getenv("SCHEDULE")
		}
		if Env.Backup == "" {
			Env.Backup = os.Getenv("BACKUP")
		}
		if Env.BackupInterval == 0 {
			Env.BackupInterval, _ = strconv.Atoi(os.Getenv("BACKUP_INTERVAL"))
		}
	}

	log.Println("env:", Env.path)
	log.Println("version:", Version, CompileDate)
	log.Println("interval:", Env.Interval)
	log.Println("schedule:", Env.Schedule)
	if Env.Interval < 1 {
		panic("interval 不可小於 1")
	}
}

func main() {

	// 建構並初始化排程
	prc = process.New(log.New(os.Stdout, "[process]", log.Lshortfile|log.LstdFlags))
	readSchedule(Env.Schedule, func(command, timeSet string, repeat int) {
		if _, err := prc.Receive(command, timeSet, repeat); err != nil {
			log.Printf("[process] ignored because Receive error: %v \n\t# %s|%s\n", err, timeSet, command)
		}
	})

	// 建構並啟動 work ticker
	work := worker.New(Env.Interval)
	go work.Run(func(now time.Time, schedule parser.SpecSchedule) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[worker run] recover %+v\n", r)
			}
		}()
		prc.Exec(now, prc.Find(schedule))
	})

	// TODO: 評估要不要拿掉
	if Env.AdminAPIAddr != "" {
		http.Handle(Env.AdminAPIPrefix+"/", http.StripPrefix(Env.AdminAPIPrefix, CreateRESTHandler()))
		go func() {
			err := http.ListenAndServe(Env.AdminAPIAddr, nil)
			if err != nil {
				log.Println("[http] ", err)
			}
		}()
		log.Println("admin api:", Env.AdminAPIAddr, Env.AdminAPIPrefix)

		// 背景備份, 僅在啟用遠端 API, 且備份間隔大於0 時才需要啟動
		if Env.BackupInterval > 0 {
			log.Println("backup:", Env.Backup, time.Duration(Env.BackupInterval)*time.Second)
			prc.Backup(Env.Backup, time.Duration(Env.BackupInterval)*time.Second)
		}
	}

	// listen os signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	signal := <-shutdown
	log.Printf("[shutdown] by %s(%#v)\n", signal, signal)

	log.Printf("system info: %+v\n", getSystemInfo())

	work.Stop()
	prc.Stop()
	prc.Wait()
}

func readSchedule(conf string, fn func(command, timeSet string, repeat int)) error {
	f, err := os.Open(conf)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line int
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
			log.Printf("line: %d schema error\n", line)
			continue
		}

		repeat, _ := strconv.Atoi(arr[1])
		command := arr[2]
		timeSet := arr[0]
		fn(command, timeSet, repeat)
		line++
	}

	return nil
}

type systemInfo struct {
	Go           string `json:"go"`
	Version      string `json:"version"`
	CompileDate  string `json:"compile-date"`
	CPU          int    `json:"CUP"`
	Goroutines   int    `json:"goroutines"`
	MemAllocated uint64 `json:"memory-allocated"`
	NextGC       string `json:"next-gc"`
}

func getSystemInfo() systemInfo {

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	return systemInfo{
		Go:           runtime.Version(),
		Version:      Version,
		CompileDate:  CompileDate,
		CPU:          runtime.NumCPU(),
		Goroutines:   runtime.NumGoroutine(),
		MemAllocated: m.Alloc,
		NextGC:       time.Duration(m.NextGC).String(),
	}
}
