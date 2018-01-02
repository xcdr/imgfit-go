package main

import (
	"fmt"
	"net/http"
	"os"
)

var srv server

func init() {

	configInit(&srv.cfg)

	if !configParse("./") {
		fmt.Println("cannot parse config!")
		os.Exit(1)
	}

	configLoad(&srv.cfg)

	http.HandleFunc("/", srv.handle)
	http.HandleFunc("/server-status", srv.status)

	if err := srv.init(); err != nil {
		fmt.Printf("cannot init server %v\n", err)
		os.Exit(1)
	}

	srv.start(nil)
}
