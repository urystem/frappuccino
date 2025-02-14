package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"

	"hot-coffee/internal/router"
)

func main() {
	port, dir := flag.String("port", "8080", "set port"), flag.String("dir", "data", "directory")
	flag.Usage = func() {
		os.Stdout.WriteString(`Coffee Shop Management System

Usage:
   hot-coffee [--port <N>] [--dir <S>]
   hot-coffee --help
	
Options:
   --help       Show this screen.
   --port N     Port number.
   --dir S      Path to the data directory.` + "\n")
	}
	flag.Parse()
	if *port = strings.TrimLeft(*port, "0 "); !regexp.MustCompile(`^\d+$`).MatchString(*port) {
		slog.Error("Invalid port")
	} else if dirinfo, err := os.Stat(*dir); err != nil {
		slog.Error("diretory error:", "", err)
	} else if !dirinfo.IsDir() {
		slog.Error("In the path was file, not directory")
	} else if err = func() error {
		for _, v := range router.PathFiles {
			if fileinfo, er := os.Stat(*dir + v); os.IsNotExist(er) {
				if e := os.WriteFile(*dir+v, []byte("[]"), 0o644); e != nil {
					return e
				}
			} else if er != nil {
				return er
			} else if fileinfo.IsDir() {
				return errors.New(v + " not file, it is a directory")
			} else if file, er := os.ReadFile(*dir + v); er != nil {
				return er
			} else if !json.Valid(file) {
				return errors.New(v + " file is not a valid json")
			}
		}
		return nil
	}(); err != nil {
		slog.Error("metadata", "error", err)
	} else {
		os.Stdout.WriteString("Dir: " + *dir + "\nPort: " + *port + "\n")
		slog.Info("Listen and Serve starting")
		slog.Error("Server stoped", "error", http.ListenAndServe(":"+*port, router.Allrouter(dir)))
	}
	os.Exit(1)
}
