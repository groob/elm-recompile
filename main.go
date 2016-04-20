package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

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
		// first time app runs
		output := compile(path)
		prev = check(output, prev)
		for {
			select {
			case <-watcher.Events:
				output := compile(path)
				prev = check(output, prev)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	err = addFoldersInPath(filepath.Dir(path), watcher)
	if err != nil {
		log.Fatal(err)
	}
	// err = watcher.Add(filepath.Dir(path))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	<-done
}

func addFoldersInPath(location string, watcher *fsnotify.Watcher) error {
	err := filepath.Walk(location, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			fmt.Println(path)
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func check(in, prev []byte) []byte {
	if !bytes.Equal(in, prev) {
		clear()
		fmt.Println(string(in))
		return in
	}
	return prev
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
