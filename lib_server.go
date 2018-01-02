package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"

	"net/url"

	"image"

	"strings"

	log "github.com/sirupsen/logrus"
)

type server struct {
	cfg        config
	watermarks map[string]image.Image
}

func (s *server) init() error {
	var err error

	s.watermarks = make(map[string]image.Image)

	for key, host := range s.cfg.Hosts {

		if len(host.Watermark) > 0 {
			s.watermarks[key], err = getWatermark(host.Watermark)

			if err != nil {
				return err
			}
		}

		if s.watermarks[key] != nil {
			log.Infof("loaded watermark (%s): %v", key, host.Watermark)
		}
	}

	return nil
}

func (s *server) start(h http.Handler) (io.Closer, error) {
	var err error
	var addr string

	var listener net.Listener

	srv := &http.Server{Addr: addr, Handler: h}

	addr = fmt.Sprintf("%s:%d", s.cfg.IP, s.cfg.Port)

	listener, err = net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	go func() {
		err := srv.Serve(listener)

		if err != nil {
			log.Printf("HTTP Server Error: %v ", err)
		}
	}()

	return io.Closer(listener), nil
}

func (s *server) getImageRequest(host string, url *url.URL) *imageRequest {
	var result imageRequest

	var ver = url.Query().Get("v")

	if matched, _ := regexp.MatchString("^[0-9]+$", ver); len(ver) > 0 && !matched {
		return nil
	}

	if len(ver) == 0 {
		ver = "default"
	}

	result.RequestPath = url.Path
	result.CachedFilePath = path.Join(s.cfg.CacheDir, host, ver, url.Path)

	signSize := rune(url.Path[1])
	result.Size = uint(s.cfg.Sizes[signSize])

	RequestPattern := "^/[1-9]/.+\\.jpg$"

	if matched, _ := regexp.MatchString(RequestPattern, url.Path); matched {
		result.SourceFilePath = path.Join(s.cfg.Hosts[host].BaseDir, url.Path[2:])
		return &result
	}

	return nil
}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	var err error

	serverName := strings.ToLower(strings.Split(r.Host, ":")[0])

	if _, ok := s.cfg.Hosts[serverName]; !ok || r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		log.Infof("404: %s", r.URL.Path)
		return
	}

	reqImage := s.getImageRequest(serverName, r.URL)

	if reqImage == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Infof("400: %s", r.URL.Path)
		return
	}

	if _, err = os.Stat(reqImage.SourceFilePath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		log.Infof("404: %s", reqImage.RequestPath)
		return
	}

	_, err = os.Stat(reqImage.CachedFilePath)

	if os.IsNotExist(err) {
		log.Infof("scale image: %s -> %s", reqImage.SourceFilePath, reqImage.CachedFilePath)

		if reqImage.resize(s.watermarks[serverName]) != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("500: cannot resize: %s", reqImage.SourceFilePath)
			return
		}
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", s.cfg.Hosts[serverName].CacheAge))

	http.ServeFile(w, r, reqImage.CachedFilePath)
	log.Infof("200: %s", reqImage.RequestPath)
	return
}

func (s *server) status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}
