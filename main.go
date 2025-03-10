package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: mylang <имя_файла.mylang>")
		return
	}

	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".mylang") {
		fmt.Println("Ошибка: файл должен иметь расширение .mylang")
		return
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	interpreter := NewInterpreter()
	interpreter.ExecuteProgram(string(content))
}