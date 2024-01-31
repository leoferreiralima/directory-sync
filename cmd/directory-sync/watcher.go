package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher           *fsnotify.Watcher
	sourceDestPairs   map[string]string
	sourceDestFolders map[string]string
}

func NewWatcher(sourceDestPairs map[string]string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		watcher:           watcher,
		sourceDestPairs:   sourceDestPairs,
		sourceDestFolders: make(map[string]string),
	}

	return w, nil
}

func (w *Watcher) WatchFolders() error {
	for source, dest := range w.sourceDestPairs {
		err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return w.watcher.Add(path)
		})
		if err != nil {
			return err
		}
		w.sourceDestFolders[source] = dest
		log.Printf("Watching folder: %s (Source) -> %s (Destination)\n", source, dest)
	}

	return nil
}

func (w *Watcher) Add(name string) error {
	return w.watcher.Add(name)
}

func (w *Watcher) Remove(name string) error {
	return w.watcher.Remove(name)
}

func (w *Watcher) Close() {
	w.watcher.Close()
}
