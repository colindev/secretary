package main

import (
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"rde-tech.vir888.com/dev/secretary/secretary.git/process"

	"github.com/ant0ine/go-json-rest/rest"
)

// CreateRESTHandler return Handler by rest
func CreateRESTHandler() http.Handler {

	app, err := rest.MakeRouter(
		rest.Post("/register", registerHandler),
		rest.Delete("/revoke", revokeHandler),
		rest.Get("/dump", dumpHandler),
		rest.Get("/running", runningCommands),
		rest.Get("/system/info", systemInfo),
	)

	if err != nil {
		log.Fatal(err)
	}

	server := rest.NewApi()
	server.Use(
		&rest.AccessLogApacheMiddleware{},
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.PoweredByMiddleware{
			XPoweredBy: "rde-tech",
		},
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
	)

	server.SetApp(app)
	return server.MakeHandler()
}

// Error for REST response
type Error struct {
	Msg  interface{} `json:"msg"`
	Code int         `json:"code"`
}

// WriteSuccess simply for use rest.ResponseWriter with struct
func WriteSuccess(w rest.ResponseWriter, ret interface{}) {
	w.WriteJson(struct {
		Success interface{} `json:"success"`
	}{ret})
}

// WriteError simply for use restResponseWriter with Error struct
func WriteError(w rest.ResponseWriter, err interface{}, code int) {
	w.WriteHeader(code)
	w.WriteJson(struct {
		Error Error `json:"error"`
	}{Error{Msg: err, Code: code}})
}

func registerHandler(w rest.ResponseWriter, r *rest.Request) {

	if err := r.ParseForm(); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	v := r.PostForm
	repeat, err := strconv.Atoi(v.Get("repeat"))
	command := v.Get("command")
	timeSet := v.Get("datetime")

	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(command, "http://") && !strings.HasPrefix(command, "https://") {
		WriteError(w, "僅接受 web hook", http.StatusBadRequest)
		return
	}

	id, err := prc.Receive(command, timeSet, repeat)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	WriteSuccess(w, id)
}

func revokeHandler(w rest.ResponseWriter, r *rest.Request) {
	prc.Revoke(r.URL.Query().Get("id"))
	WriteSuccess(w, "ok")
}

func dumpHandler(w rest.ResponseWriter, r *rest.Request) {

	var ret = map[string]string{}
	prc.Each(func(c *process.Command) error {
		ret[c.ID] = c.String()
		return nil
	})
	WriteSuccess(w, ret)
}

func runningCommands(w rest.ResponseWriter, r *rest.Request) {
	WriteSuccess(w, prc.Running())
}

func systemInfo(w rest.ResponseWriter, r *rest.Request) {

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	WriteSuccess(w, struct {
		Go           string `json:"go"`
		Version      string `json:"version"`
		CompileDate  string `json:"compile-date"`
		CPU          int    `json:"CUP"`
		Goroutines   int    `json:"goroutines"`
		MemAllocated uint64 `json:"memory-allocated"`
		NextGC       string `json:"next-gc"`
	}{
		Go:           runtime.Version(),
		Version:      Version,
		CompileDate:  CompileDate,
		CPU:          runtime.NumCPU(),
		Goroutines:   runtime.NumGoroutine(),
		MemAllocated: m.Alloc,
		NextGC:       time.Duration(m.NextGC).String(),
	})
}
