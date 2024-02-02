package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileOps struct {
	sourceFolder string
	destFolder   string
	watcher      *Watcher
}

func NewFileOps(sourceFolder, destFolder string, watcher *Watcher) *FileOps {
	return &FileOps{
		sourceFolder: sourceFolder,
		destFolder:   destFolder,
		watcher:      watcher,
	}
}

func (f *FileOps) CopyFile(src string) {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Println("Error opening source file:", err)
		return
	}
	defer srcFile.Close()

	destFile := filepath.Join(f.destFolder, f.GetFilePath(src))
	dest, err := os.Create(destFile)
	if err != nil {
		log.Println("Error creating destination file:", err)
		return
	}
	defer dest.Close()

	_, err = io.Copy(dest, srcFile)
	if err != nil {
		log.Println("Error copying file:", err)
		return
	}

	log.Printf("Copied: %s -> %s\n", src, destFile)
}

func (f *FileOps) DeleteFile(src string) {
	destFile := filepath.Join(f.destFolder, f.GetFilePath(src))
	err := os.Remove(destFile)
	if err != nil {
		log.Println("Error deleting file:", err)
		return
	}

	log.Printf("Deleted: %s\n", destFile)
}

func (f *FileOps) IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println("Error checking if directory:", err)
		return false
	}
	return fileInfo.IsDir()
}

func (f *FileOps) CreateFolder(src string) {
	dest := filepath.Join(f.destFolder, f.GetFilePath(src))
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		log.Println("Error creating folder:", err)
		return
	}

	f.watcher.Add(src)
	log.Printf("Created folder: %s -> %s\n", src, dest)
}

func (f *FileOps) DeleteFolder(src string) {
	dest := filepath.Join(f.destFolder, f.GetFilePath(src))
	err := os.RemoveAll(dest)
	if err != nil {
		log.Println("Error deleting folder:", err)
		return
	}

	f.watcher.Remove(src)
	log.Printf("Deleted folder: %s -> %s\n", src, dest)
}

func (f *FileOps) GetFilePath(src string) string {
	return strings.Replace(src, strings.Replace(f.sourceFolder, ".\\", "", 1), "", 1)
}
