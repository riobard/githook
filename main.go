package main

import (
	"encoding/json"
	"flag"
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

	f, err := os.Open(arg.conf)
	if err != nil {
		log.Fatalf("failed to open config file %s: %s", arg.conf, err)
	}

	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		log.Fatalf("failed to parse config file %s: %s", arg.conf, err)
	}

	if arg.version {
		println(Version)
		return
	}

	log.Printf("Listening on %s", arg.addr)
	if err := http.ListenAndServe(arg.addr, http.HandlerFunc(handle)); err != nil {
		log.Fatal(err)
	}
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
