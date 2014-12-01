package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/riobard/githook/github"
)

var Version string
var conf map[string]*Hook

var arg struct {
	version bool
	conf    string
	addr    string
}

type Hook struct {
	Source  string
	Secret  string
	Command string
}

func main() {
	flag.StringVar(&arg.addr, "addr", ":4008", "listening address")
	flag.StringVar(&arg.conf, "conf", "/etc/githook.conf", "path to config file")
	flag.BoolVar(&arg.version, "version", false, "print version number")
	flag.Parse()

	if arg.version {
		println(Version)
		return
	}

	f, err := os.Open(arg.conf)
	if err != nil {
		log.Fatalf("failed to open config file %s: %s", arg.conf, err)
	}

	conf, err = parseConf(f)
	if err != nil {
		log.Fatalf("failed to parse config file %s: %s", arg.conf, err)
	}

	log.Printf("Listening on %s", arg.addr)
	if err := http.ListenAndServe(arg.addr, http.HandlerFunc(handle)); err != nil {
		log.Fatalf("http.ListenAndServe error: %s", err)
	}
}

func parseConf(r io.Reader) (map[string]*Hook, error) {
	conf := make(map[string]*Hook)
	if err := json.NewDecoder(r).Decode(&conf); err != nil {
		return nil, err
	}
	return conf, nil

}

func handle(w http.ResponseWriter, r *http.Request) {
	hook, ok := conf[r.URL.Path]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	switch src := strings.ToLower(hook.Source); src {
	case "github":
		github.NewHook(hook.Secret, hook.Command).ServeHTTP(w, r)
	default:
		http.Error(w, "unexpected source", http.StatusBadRequest)
	}
}
