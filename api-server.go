package main

import (
	"net/http"
)

func register(rw http.ResponseWriter, req *http.Request) {
	// TODO: err to json {"error": "xxxx"}
	if err := req.ParseForm(); err != nil {
		rw.Write([]byte("error"))
		return
	}

	err := processes.Recieve(req.PostForm)

	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write([]byte("ok"))
}

func ListenAndServe(addr string) <-chan error {

	err := make(chan error)
	http.HandleFunc("/api/register", register)

	go func() {
		err <- http.ListenAndServe(addr, nil)
	}()

	// Process, http fatal
	return err
}
