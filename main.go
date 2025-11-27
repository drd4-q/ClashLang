package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: pylang <имя_файла.pyl>")
		fmt.Println("\nПример:")
		fmt.Println("  pylang example.pyl")
		return
	}

	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".pyl") {
		fmt.Println("Ошибка: файл должен иметь расширение .pyl")
		return
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	// Создаем лексер
	lexer := NewLexer(string(content))

	// Создаем парсер
	parser := NewParser(lexer)

	// Парсим программу
	program, err := parser.Parse()
	if err != nil {
		fmt.Printf("Ошибка парсинга: %v\n", err)
		return
	}

	// Создаем интерпретатор
	interpreter := NewNewInterpreter()

	// Выполняем программу
	if err := interpreter.Execute(program); err != nil {
		fmt.Printf("Ошибка выполнения: %v\n", err)
		return
	}
}
