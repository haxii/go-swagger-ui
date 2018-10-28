package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/haxii/daemon"
)

var (
	port = flag.Int("p", 8080, "server port")
	_    = flag.String("s", daemon.UsageDefaultName, daemon.UsageMessage)
)

func main() {
	daemon.Make("-s", "httpdaemon", "simple http daemon service").Run(serve)
}

func serve() {
	flag.Parse()
	http.HandleFunc("/hello",
		func(w http.ResponseWriter, req *http.Request) {
			io.WriteString(w, "hello, world!\n")
		},
	)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
