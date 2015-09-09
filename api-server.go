package main

import (
	"net/http"
)

func ListenAndServe(addr string) (<-chan *Process, <-chan error) {

	cp := make(chan *Process)
	ce := make(chan error)

	http.HandleFunc("/api/register", func(rw http.ResponseWriter, req *http.Request) {

		// TODO: err to json {"error": "xxxx"}
		if err := req.ParseForm(); err != nil {
			rw.Write([]byte("error"))
			return
		}

		p, err_p := ProcessesParse(req.PostForm)

		if err_p != nil {
			rw.Write([]byte("error"))
			return
		}

		cp <- p

		rw.Write([]byte("ok"))
	})

	go func() {
		ce <- http.ListenAndServe(addr, nil)
	}()

	// Process, http fatal
	return cp, ce
}
