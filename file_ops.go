package main

import (
	"fmt"
	"os"
)

type FileManager struct {
	currentFile string
	content     string
}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (fm *FileManager) Open(filename string) {
	fm.currentFile = filename
	if data, err := os.ReadFile(filename); err == nil {
		fm.content = string(data)
	}
}

func (fm *FileManager) Close() {
	fm.currentFile = ""
}

func (fm *FileManager) Write(text string) {
	if fm.currentFile == "" {
		fmt.Println("Error: no file open")
		return
	}
	os.WriteFile(fm.currentFile, []byte(text), 0644)
}

func (fm *FileManager) Read() {
	if fm.currentFile == "" {
		fmt.Println("Error: no file open")
		return
	}
	if data, err := os.ReadFile(fm.currentFile); err == nil {
		fm.content = string(data)
	} else {
		fmt.Printf("Error reading file: %v\n", err)
	}
}

func (fm *FileManager) Create(filename string) {
	os.WriteFile(filename, []byte{}, 0644)
}

func (fm *FileManager) GetContent() string {
	return fm.content
}

func (fm *FileManager) IsOpen() bool {
	return fm.currentFile != ""
}
