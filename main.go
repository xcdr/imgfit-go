package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var (
	version string
	build   string
)

func main() {
	var err error
	var srv server
	var cfgDir string

	log.Infof("Starting %s version: %s+%s", path.Base(os.Args[0]), version, build)

	flag.StringVar(&cfgDir, "cfg-dir", "/opt/imgfit/etc", "Config dir path")
	flag.Parse()

	configInit(&srv.cfg)
	if !configParse(cfgDir) {
		os.Exit(1)
	}

	configLoad(&srv.cfg)

	if err = srv.init(); err != nil {
		log.Fatalf("cannot load watermark %v", err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Infof("Handled signal: %s", sig)
		done <- true
	}()

	log.Infof("Start server on %s:%d", srv.cfg.IP, srv.cfg.Port)

	http.HandleFunc("/", srv.handle)
	http.HandleFunc("/server-status", srv.status)

	srvCloser, err := srv.start(nil)
	if err != nil {
		log.Fatalf("Error on server.start: %s", err)
	}

	// Waiting for signal
	<-done

	// Close HTTP Server
	err = srvCloser.Close()

	if err != nil {
		log.Fatalf("Error on server close %v", err)
	}

	log.Infoln("Stopping server")

	os.Exit(0)
}
