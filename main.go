package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: clashlang <file.clash> [--debug]")
		return
	}

	filename := os.Args[1]
	debugMode := false

	if len(os.Args) > 2 && os.Args[2] == "--debug" {
		debugMode = true
	}

	if !strings.HasSuffix(filename, ".clash") {
		fmt.Println("Error: file must have .clash extension")
		return
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	interpreter := NewInterpreter()

	if debugMode {
		fmt.Println("=== COMPILE DEBUG ===")
		fmt.Printf("File: %s\n", filename)
		fmt.Printf("Lines: %d\n", len(strings.Split(string(content), "\n")))
		fmt.Println("====================")
		fmt.Println()
	}

	interpreter.ExecuteProgram(string(content))

	if debugMode {
		fmt.Println()
		fmt.Println("=== POST-EXECUTE DEBUG ===")
		interpreter.memory.PrintDebug()
	}
}
