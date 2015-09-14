package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	addr := flag.String("addr", "", "api listen")
	interval := flag.Int64("interval", 10, "interval (seconds)")

	flag.Parse()

	worker.Run(*interval)

	ce := ListenAndServe(*addr)

	for {
		select {
		case e := <-ce:
			fmt.Println("server down", e)
			os.Exit(1)
		}
	}

}
