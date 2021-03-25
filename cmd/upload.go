package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	ssh2 "github.com/deliangyang/sftp-upload/internal/ssh"
	"github.com/fsnotify/fsnotify"
	"github.com/json-iterator/go"
)

var (
	conf string
)

func init() {
	flag.StringVar(&conf, "conf", "", "conf filename")
	flag.Parse()
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	f, _ := os.Open(conf)

	// Close the file after it has been copied
	content, _ := ioutil.ReadAll(f)
	f.Close()

	var u ssh2.User
	if err := jsoniter.Unmarshal(content, &u); err != nil {
		panic(err)
	}

	path := u.Watch
	client := ssh2.NewClient(path, u)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {

				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if strings.HasSuffix(event.Name, "~") {
						continue
					}
					log.Println("modified file:", event.Name)
					client.Upload(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done

	client.Close()
}
