package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"gopkg.in/fsnotify.v1"
)

var (
	version = "unknown"
	gitHash = "unknown"
)

func main() {
	var (
		flPort    = flag.String("port", "3000", "http port to listen to, defaults to 3000")
		flVersion = flag.Bool("version", false, "print version information")
	)
	flag.Parse()
	if *flVersion {
		fmt.Printf("elm-recompile - Version %s\n", version)
		fmt.Printf("Git Hash - %s\n", gitHash)
		os.Exit(0)
	}
	go serve(*flPort)
	inFile := os.Args[1]
	reCompile(inFile)
}

func serve(port string) {
	port = ":" + port
	log.Fatal(http.ListenAndServe(port, http.FileServer(http.Dir("."))))
}

func reCompile(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		prev := []byte("")
		for {
			select {
			case <-watcher.Events:
				output := compile(path)
				if !bytes.Equal(output, prev) {
					clear()
					fmt.Println(string(output))
					prev = output
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func compile(path string) []byte {
	makeCmd := exec.Command("/usr/local/bin/elm-make", path, "--output", "elm.js")
	output, _ := makeCmd.CombinedOutput()
	return output
}

func clear() {
	clearCmd := exec.Command("/usr/bin/clear")
	clearCmd.Stdout = os.Stdout
	err := clearCmd.Run()
	if err != nil {
		log.Fatal(err)
	}

}
