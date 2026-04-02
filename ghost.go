package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const HELP = `usage: ghost [--ip <ip or domain>] [--port <port>] [<directory>]

<directory> | The directory of the webpage (default: $PWD)
OPTIONS:
	--ip <ip or domain> | The ip or domain listened on (default: localhost)
	--port <port> | The port listened on (default: auto)
	--help | Show this text`

const PREFIX = "\x1b[48;2;191;31;127mGhost:\x1b[0m"
const GREY = "\x1b[38;2;127;127;127m"

func main() {
	log.SetPrefix(PREFIX + " \x1b[48;2;223;0;0merror:\x1b[0m ")
	flag_ip := flag.String("ip", "localhost", "--ip <ip or domain> (default: 'localhost')")
	flag_port := flag.String("port", "auto", "--port <port> (default: 80)")
	flag_help := flag.Bool("help", false, "--help")
	flag.Parse()
	if *flag_help {
		fmt.Println(HELP)
		os.Exit(0)
	}
	dir := flag.Arg(0)
	if dir == "" {
		dir = "."
	}

	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}

	if isDir("errors") {
		e404, err := os.ReadFile("errors/404.html")
		if err == nil {
			PAGERR404 = string(e404)
		}
	}

	validp := []string{"/"}

	if !isFile("./index.html") {
		log.Fatal("Main directory must contain index.html")
	}

	err = filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if strings.HasPrefix(d.Name(), "secret") {
			return nil
		}
		rel, err := filepath.Rel(".", path)
		p := "/" + rel
		if rel == "." {
			return nil
		}
		if err != nil {
			return err
		}

		if d.IsDir() && isFile(path+"/index.html") {
			data, err := os.ReadFile(path + "/index.html")
			if err != nil {
				return err
			}
			http.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				_, _ = w.Write(data)
			})
			validp = append(validp, p)
		} else if !d.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var contype string = http.DetectContentType(data)
			if strings.HasSuffix(path, ".css") {
				contype = "text/css"
			} else if strings.HasSuffix(path, ".js") {
				contype = "application/javascript"
			}
			http.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", contype)
				_, _ = w.Write(data)
			})
			validp = append(validp, p)
		}
		return nil
	})
	if err != nil {
		fmt.Println(PREFIX, "ERROR:", err)
	}

	idx, err := os.ReadFile("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(validp, r.URL.Path) {
			_, _ = w.Write([]byte(PAGERR404))
			return
		}
		_, _ = w.Write(idx)
	})
	addr := *flag_ip + ":" + *flag_port

	if isFile("secret/cert.pem") && isFile("secret/key.pem") {
		addr = strings.ReplaceAll(addr, ":auto", ":443")
		fmt.Println(PREFIX, "Starting listening on "+addr+"...")
		err := http.ListenAndServeTLS(addr, "secret/cert.pem", "secret/key.pem", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(PREFIX + GREY + "Not using TLS due to lack of a certificate!\x1b[0m")
		addr = strings.ReplaceAll(addr, ":auto", ":80")
		fmt.Println(PREFIX, "Starting listening on "+addr+"...")
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func isFile(path string) bool {
	i, err := os.Stat(path)
	if err != nil {
		return false
	}
	if i.IsDir() {
		return false
	}
	return true
}

func isDir(path string) bool {
	i, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !i.IsDir() {
		return false
	}
	return true
}
