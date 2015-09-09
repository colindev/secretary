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

	receive := ProcessesReceiver(*interval)

	cp, ce := ListenAndServe(*addr)

	for {
		select {
		case p := <-cp:
			if e := receive(p); e != nil {
				fmt.Println("reject", e)
			} else {
				fmt.Println("recieve", p)
			}
			break

		case e := <-ce:
			fmt.Println("server down", e)
			os.Exit(1)
		}
	}

}
