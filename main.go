package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/fsnotify.v1"
)

var (
	version = "unknown"
	gitHash = "unknown"
)

func main() {
	inFile := os.Args[1]
	reCompile(inFile)
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
