package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Command описывает структуру команды из JSON
type Command struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Pattern     string `json:"pattern"`
}

// CommandList для десериализации JSON
type CommandList struct {
	Commands []Command `json:"commands"`
}

// Interpreter управляет выполнением программы
type Interpreter struct {
	commands   map[int]Command
	variables  map[string]interface{}
	lastResult interface{}
	functions  map[string][]string
	inIfBlock  bool
	ifBlock    []string
}

func NewInterpreter() *Interpreter {
	interp := &Interpreter{
		commands:  make(map[int]Command),
		variables: make(map[string]interface{}),
		functions: make(map[string][]string),
		ifBlock:   []string{},
	}
	interp.loadCommands()
	return interp
}

func (i *Interpreter) loadCommands() {
	file, err := os.ReadFile("commands.json")
	if err != nil {
		fmt.Println("Ошибка загрузки commands.json:", err)
		return
	}
	var cmdList CommandList
	if err := json.Unmarshal(file, &cmdList); err != nil {
		fmt.Println("Ошибка разбора JSON:", err)
		return
	}
	for _, cmd := range cmdList.Commands {
		i.commands[cmd.ID] = cmd
	}
}

func (i *Interpreter) matchCommand(line string) (Command, map[string]string, bool) {
	line = strings.TrimSpace(line)
	for _, cmd := range i.commands {
		pattern := cmd.Pattern
		parts := strings.Split(pattern, "{{")
		if len(parts) == 1 && line == pattern {
			return cmd, nil, true
		}
		if strings.HasPrefix(line, parts[0]) {
			remainder := strings.TrimPrefix(line, parts[0])
			if len(parts) > 1 {
				params := make(map[string]string)
				for _, part := range parts[1:] {
					end := strings.Index(part, "}}")
					if end == -1 {
						continue
					}
					varName := part[:end]
					closing := part[end+2:]
					if strings.HasSuffix(remainder, closing) {
						value := strings.TrimSuffix(remainder, closing)
						params[varName] = value
						return cmd, params, true
					} else {
						nextPartIdx := strings.Index(remainder, " ")
						if nextPartIdx != -1 {
							params[varName] = remainder[:nextPartIdx]
							remainder = strings.TrimSpace(remainder[nextPartIdx:])
						}
					}
				}
			}
		}
	}
	return Command{}, nil, false
}

func (i *Interpreter) ExecuteStatement(line string) {
	// Обрезаем комментарии: всё после "//"
	if idx := strings.Index(line, "//"); idx != -1 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	switch {
	case strings.HasPrefix(line, "Memory load ("):
		funcNames := strings.Trim(strings.TrimPrefix(line, "Memory load ("), ")")
		for _, funcName := range strings.Split(funcNames, ",") {
			funcName = strings.TrimSpace(funcName)
			if cmds, ok := i.functions[funcName]; ok {
				for _, cmd := range cmds {
					i.ExecuteStatement(cmd)
				}
			}
		}
	case line == "Memory out":
		fmt.Println("Global memory:", i.variables)
	case line == "}":
		if i.inIfBlock {
			i.inIfBlock = false
			for _, cmd := range i.ifBlock {
				i.ExecuteStatement(cmd)
			}
			i.ifBlock = []string{}
		}
	default:
		cmd, params, matched := i.matchCommand(line)
		if !matched {
			if i.inIfBlock {
				i.ifBlock = append(i.ifBlock, line)
			}
			return
		}
		switch cmd.ID {
		case 1: // Print
			varName := params["var"]
			if val, ok := i.variables[varName]; ok {
				fmt.Println(val)
			} else {
				fmt.Println(varName) // Вывод строки, если переменная не определена
			}
		case 2: // Solve.input
			varName := params["var"]
			fmt.Printf("Введите число для %s: ", varName)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			num, _ := strconv.Atoi(strings.TrimSpace(input))
			i.variables[varName] = num
		case 3: // Solve
			expr := params["expr"]
			i.lastResult = parseExpression(expr, i.variables)
		case 4: // Solve.out
			varName := params["var"]
			i.variables[varName] = i.lastResult
		case 5: // Text.input
			varName := params["var"]
			fmt.Printf("Введите текст для %s: ", varName)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			i.variables[varName] = strings.TrimSpace(input)
		case 6: // Text
			expr := params["expr"]
			i.lastResult = parseTextExpression(expr, i.variables)
		case 7: // Text.out
			varName := params["var"]
			i.variables[varName] = i.lastResult
		case 8: // If
			varName := params["var"]
			valueStr := params["value"]
			value, _ := strconv.Atoi(valueStr)
			if val, ok := i.variables[varName].(int); ok && val == value {
				i.inIfBlock = true
			}
		case 9: // jump
			funcName := params["func"]
			if cmds, ok := i.functions[funcName]; ok {
				for _, cmd := range cmds {
					i.ExecuteStatement(cmd)
				}
			}
		case 13: // Text.length
			varName := params["var"]
			if val, ok := i.variables[varName].(string); ok {
				i.lastResult = utf8.RuneCountInString(val)
			} else {
				fmt.Println("Ошибка: переменная не является текстом")
			}
		case 14: // Text.upper
			varName := params["var"]
			if val, ok := i.variables[varName].(string); ok {
				i.lastResult = strings.ToUpper(val)
			} else {
				fmt.Println("Ошибка: переменная не является текстом")
			}
		}
	}
}

func parseExpression(expr string, variables map[string]interface{}) int {
	expr = strings.ReplaceAll(expr, " ", "")
	parts := strings.Split(expr, "+")
	if len(parts) == 2 {
		left := getValue(parts[0], variables)
		right := getValue(parts[1], variables)
		return left + right
	}
	return getValue(expr, variables)
}

func getValue(part string, variables map[string]interface{}) int {
	if val, ok := variables[part]; ok {
		return val.(int)
	}
	num, _ := strconv.Atoi(part)
	return num
}

func parseTextExpression(expr string, variables map[string]interface{}) string {
	expr = strings.ReplaceAll(expr, " ", "")
	parts := strings.Split(expr, "+")
	var result string
	for _, part := range parts {
		if val, ok := variables[part].(string); ok {
			if result != "" {
				result += " "
			}
			result += val
		} else {
			fmt.Println("Ошибка: переменная не является текстом")
		}
	}
	return result
}

func (i *Interpreter) ExecuteProgram(program string) {
	lines := strings.Split(program, "\n")
	var currentFunction string

	for _, line := range lines {
		// Обрезаем комментарии: всё после "//"
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Function ("):
			currentFunction = strings.Trim(strings.TrimPrefix(line, "Function ("), ")")
			i.functions[currentFunction] = []string{}
		case line == "Memory start (":
			continue
		case line == ")":
			currentFunction = ""
		case currentFunction != "":
			i.functions[currentFunction] = append(i.functions[currentFunction], line)
		default:
			i.ExecuteStatement(line)
		}
	}
}