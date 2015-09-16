package main

import (
	"errors"
	"fmt"
	"net/http"
)

type Handle func(http.ResponseWriter, *http.Request, *Process)

type Router struct {
	process *Process
	routes  map[string]map[string]Handle
}

func (r *Router) On(method string, path string, handler Handle) {
	r.routes[method][path] = handler
}

func (r *Router) Get(path string, handler Handle) {
	r.On("GET", path, handler)
}

func (r *Router) Post(path string, handler Handle) {
	r.On("POST", path, handler)
}

func (r *Router) Put(path string, handler Handle) {
	r.On("PUT", path, handler)
}

func (r *Router) Delete(path string, handler Handle) {
	r.On("DELETE", path, handler)
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	handlers, ok := r.routes[req.Method]

	if !ok {
		http.NotFound(rw, req)
		return
	}

	handler, ok := handlers[req.URL.Path]
	if !ok {
		http.NotFound(rw, req)
		return
	}

	handler(rw, req, r.process)

}

func NewRouter(p *Process) *Router {
	r := &Router{}
	r.process = p
	r.routes = make(map[string]map[string]Handle)

	return r
}

func ListenAndServe(addr string, p *Process) <-chan error {

	err := make(chan error)
	router := NewRouter(p)
	router.Post("/api/register", registerHandler)
	router.Delete("/api/revoke", revokeHandler)
	router.Get("/api/dump", dumpHandler)

	go func() {
		err <- http.ListenAndServe(addr, router)
	}()

	// Process, http fatal
	return err
}

func jsonReply(rw http.ResponseWriter, msg string) {
	h := rw.Header()
	h.Add("Content-Type", "application/json")
	rw.Write([]byte(msg))
}

func jsonReplyError(rw http.ResponseWriter, e error) {
	jsonReply(rw, fmt.Sprintf("{\"error\":\"%s\"}", e.Error()))
}

func jsonReplySuccess(rw http.ResponseWriter) {
	jsonReply(rw, "{\"ok\":true}")
}

func registerHandler(rw http.ResponseWriter, req *http.Request, p *Process) {

	if "POST" != req.Method {
		jsonReplyError(rw, errors.New("請使用post"))
	}

	if err := req.ParseForm(); err != nil {
		jsonReplyError(rw, err)
		return
	}

	if err := p.Recieve(req.PostForm); err != nil {
		jsonReplyError(rw, err)
		return
	}

	jsonReplySuccess(rw)
}

func revokeHandler(rw http.ResponseWriter, req *http.Request, p *Process) {
	if err := req.ParseForm(); err != nil {
		jsonReplyError(rw, err)
		return
	}

	p.Revoke(req.Form.Get("id"))
}

func dumpHandler(rw http.ResponseWriter, req *http.Request, p *Process) {
	rw.Write([]byte(p.Dump()))
}
