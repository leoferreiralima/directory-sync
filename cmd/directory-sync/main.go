package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if len(os.Args) < 3 || len(os.Args)%2 != 1 {
		log.Fatal("Usage: go run main.go source1 dest1 source2 dest2 ...")
	}

	sourceDestPairs := make(map[string]string)
	for i := 1; i < len(os.Args)-1; i += 2 {
		sourceDestPairs[os.Args[i]] = os.Args[i+1]
	}

	w, err := NewWatcher(sourceDestPairs)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	if err := w.WatchFolders(); err != nil {
		log.Fatal(err)
	}

	ops := make(map[string]*FileOps)
	for source, dest := range w.sourceDestFolders {
		ops[source] = NewFileOps(source, dest, w)
	}

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			sourceFolder := getSourceFolder(event.Name, w.sourceDestPairs)
			destFolder := w.sourceDestFolders[sourceFolder]

			log.Printf("Event.Name: %s\n", event.Name)

			ops[sourceFolder].processEvent(event)

			log.Printf("Event (%s): %s (Source) -> %s (Destination)\n", event.Op, event.Name, destFolder)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

func (f *FileOps) processEvent(event fsnotify.Event) {
	switch event.Op {
	case fsnotify.Create, fsnotify.Write, fsnotify.Rename:
		if f.isDirectory(event.Name) {
			f.createFolder(event.Name)
		} else {
			f.CopyFile(event.Name)
		}
	case fsnotify.Remove:
		if f.isDirectory(event.Name) {
			f.deleteFolder(event.Name)

		} else {
			f.DeleteFile(event.Name)
		}
	}
}

func getSourceFolder(filePath string, sourceDestPairs map[string]string) string {
	for source := range sourceDestPairs {
		if isSubPath(source, filePath) {
			return source
		}
	}
	return ""
}

func isSubPath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
