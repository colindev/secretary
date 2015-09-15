package main

import (
	"net/http"
)

func ListenAndServe(addr string, p *Process) <-chan error {

	err := make(chan error)
	http.HandleFunc("/api/register", getRegister(p))

	go func() {
		err <- http.ListenAndServe(addr, nil)
	}()

	// Process, http fatal
	return err
}

func getRegister(p *Process) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		// TODO: err to json {"error": "xxxx"}
		if err := req.ParseForm(); err != nil {
			rw.Write([]byte("error"))
			return
		}

		err := p.Recieve(req.PostForm)

		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Write([]byte("ok"))
	}
}
