package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: clashlang <имя_файла.clash>")
		return
	}

	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".clash") {
		fmt.Println("Ошибка: файл должен иметь расширение .clash")
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