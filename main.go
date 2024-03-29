package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var from string

type publisher struct {
	*fsnotify.Watcher
	subscribers map[subscriber]bool
}
func (p *publisher) subscribe(s subscriber) {
	p.subscribers[s] = true
}
func (p *publisher) unsubscribe(s subscriber) {
	delete(p.subscribers, s)
}
func (p *publisher) notify(e fsnotify.Event) {
	for sub, ok := range p.subscribers {
		if !ok {
			delete(p.subscribers, sub)
			return
		}
		sub.update(e)
	}
}
func (p *publisher) serve() {
	for {
		select {
		case event := <-p.Events:
			p.notify(event)
		case err := <-p.Errors:
			if err != nil {
				log.Printf("Watcher Error: %v\n", err)
				return
			}
		}
	}
}
func createPublisher() *publisher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Cannot create watcher: %v\n", err)
	}
	return &publisher{
		Watcher:     watcher,
		subscribers: make(map[subscriber]bool),
	}
}

type subscriber interface {
	update(event fsnotify.Event)
}
/* service implements the Subscriber Interface */
type service struct {
	from string
}
func (s service) update(e fsnotify.Event) {
	r := regexp.MustCompile(`.*/\..*`)
	isDotFile := r.MatchString(e.Name)

	if e.Op != fsnotify.Create || isDotFile {
		return
	}
	// the file in question was added to the folder and is not a dot file

	files, err := ioutil.ReadDir(s.from)
	if err != nil {
		log.Fatalf("could not read directory %s: %v\n", s.from, err)
	}

	if len(files) <= 10 {
		return
	}

	sort.Slice(files, func(i, j int) bool {
		fileA, fileB := files[i], files[j]
		return fileA.ModTime().After(fileB.ModTime())
	})

	cutoff := files[9].ModTime()

	err = filepath.Walk(s.from, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		modified := info.ModTime()
		if modified.Before(cutoff) {
			// in other words, the modified time of the file is older than the 10th newest
			err = macOSTrash(path)
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error walking the folder: %v", err)
	}
}
func createService(from string) *service {
	return &service{
		from: from,
	}
}

// macOSTrash moves a file or folder including its content into the systems trashbin.
// The path MUST be an absolute path otherwise the executable will not run
func macOSTrash(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	// catch any unknown error
	if err != nil {
		return err
	}

	path = strings.ReplaceAll(path, "\"", "\\\"")
	osascriptCommand := fmt.Sprintf("tell app \"Finder\" to delete POSIX file \"%s\"", path)
	return exec.Command("osascript", "-e", osascriptCommand).Run()
}

func init() {
	flag.StringVar(&from, "from", "/Invalid/Path", "the absolute path of the screenshot folder")
}

func main() {
	flag.Parse()
	if _, err := os.Stat(from); err != nil {
		log.Fatalf("from path is not reachable. Please pass a valid path to your screenshot directory as the -from flag: %v\n", err)
	}
	publisher := createPublisher()
	defer publisher.Close()
	service := createService(from)
	publisher.subscribe(service)

	go publisher.serve()
	publisher.Add(from)
	log.Println("Watching for changes from", from, "...")
	<-make(chan bool)
}
