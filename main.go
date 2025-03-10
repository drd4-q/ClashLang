package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: clashlang <имя_файла.clashlang>")
		return
	}

	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".clashlang") {
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