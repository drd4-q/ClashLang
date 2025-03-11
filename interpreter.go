package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел
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
	lineLower := strings.ToLower(line)

	for _, cmd := range i.commands {
		pattern := strings.ToLower(cmd.Pattern)
		parts := strings.Split(pattern, "{{")
		if len(parts) == 1 && lineLower == pattern {
			return cmd, nil, true
		}
		if strings.HasPrefix(lineLower, parts[0]) {
			remainder := strings.TrimPrefix(line, cmd.Pattern[:len(parts[0])])
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
	if idx := strings.Index(line, "//"); idx != -1 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	lineLower := strings.ToLower(line)

	switch {
	case strings.HasPrefix(lineLower, "memory load ("):
		funcNames := strings.Trim(strings.TrimPrefix(line, "Memory load ("), ")")
		for _, funcName := range strings.Split(funcNames, ",") {
			funcName = strings.TrimSpace(funcName)
			if cmds, ok := i.functions[funcName]; ok {
				for _, cmd := range cmds {
					i.ExecuteStatement(cmd)
				}
			}
		}
	case lineLower == "}":
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
				fmt.Println(varName)
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
		case 10: // memory out
			// Форматированный вывод переменных
			var keys []string
			for k := range i.variables {
				keys = append(keys, k)
			}
			sort.Strings(keys) // Сортировка ключей по алфавиту
			for idx, key := range keys {
				if val, ok := i.variables[key]; ok {
					fmt.Printf("%d) %v\n", idx+1, val)
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
		// Новые функции
		case 15: // abs
			varName := params["var"]
			if val, ok := i.variables[varName].(float64); ok {
				i.lastResult = math.Abs(val)
			} else if val, ok := i.variables[varName].(int); ok {
				i.lastResult = float64(math.Abs(float64(val)))
			} else {
				fmt.Println("Ошибка: переменная не является числом")
			}
		case 16: // sqrt
			varName := params["var"]
			if val, ok := i.variables[varName].(float64); ok {
				i.lastResult = math.Sqrt(val)
			} else if val, ok := i.variables[varName].(int); ok {
				i.lastResult = math.Sqrt(float64(val))
			} else {
				fmt.Println("Ошибка: переменная не является числом")
			}
		case 17: // pow
			baseStr := params["base"]
			expStr := params["exponent"]
			var base, exp float64
			if val, ok := i.variables[baseStr]; ok {
				if v, ok := val.(int); ok {
					base = float64(v)
				} else if v, ok := val.(float64); ok {
					base = v
				}
			} else {
				base, _ = strconv.ParseFloat(baseStr, 64)
			}
			if val, ok := i.variables[expStr]; ok {
				if v, ok := val.(int); ok {
					exp = float64(v)
				} else if v, ok := val.(float64); ok {
					exp = v
				}
			} else {
				exp, _ = strconv.ParseFloat(expStr, 64)
			}
			i.lastResult = math.Pow(base, exp)
		case 18: // round
			varName := params["var"]
			if val, ok := i.variables[varName].(float64); ok {
				i.lastResult = math.Round(val)
			} else if val, ok := i.variables[varName].(int); ok {
				i.lastResult = float64(val)
			} else {
				fmt.Println("Ошибка: переменная не является числом")
			}
		case 19: // sin
			varName := params["var"]
			if val, ok := i.variables[varName].(float64); ok {
				i.lastResult = math.Sin(val)
			} else if val, ok := i.variables[varName].(int); ok {
				i.lastResult = math.Sin(float64(val))
			} else {
				fmt.Println("Ошибка: переменная не является числом")
			}
		case 20: // cos
			varName := params["var"]
			if val, ok := i.variables[varName].(float64); ok {
				i.lastResult = math.Cos(val)
			} else if val, ok := i.variables[varName].(int); ok {
				i.lastResult = math.Cos(float64(val))
			} else {
				fmt.Println("Ошибка: переменная не является числом")
			}
		case 21: // tan
			varName := params["var"]
			if val, ok := i.variables[varName].(float64); ok {
				i.lastResult = math.Tan(val)
			} else if val, ok := i.variables[varName].(int); ok {
				i.lastResult = math.Tan(float64(val))
			} else {
				fmt.Println("Ошибка: переменная не является числом")
			}
		case 22: // log
			varName := params["var"]
			if val, ok := i.variables[varName].(float64); ok {
				i.lastResult = math.Log(val)
			} else if val, ok := i.variables[varName].(int); ok {
				i.lastResult = math.Log(float64(val))
			} else {
				fmt.Println("Ошибка: переменная не является числом")
			}
		case 23: // log10
			varName := params["var"]
			if val, ok := i.variables[varName].(float64); ok {
				i.lastResult = math.Log10(val)
			} else if val, ok := i.variables[varName].(int); ok {
				i.lastResult = math.Log10(float64(val))
			} else {
				fmt.Println("Ошибка: переменная не является числом")
			}
		case 24: // random
			i.lastResult = rand.Float64()
		case 25: // randint
			minStr := params["min"]
			maxStr := params["max"]
			var min, max int
			if val, ok := i.variables[minStr]; ok {
				if v, ok := val.(int); ok {
					min = v
				}
			} else {
				min, _ = strconv.Atoi(minStr)
			}
			if val, ok := i.variables[maxStr]; ok {
				if v, ok := val.(int); ok {
					max = v
				}
			} else {
				max, _ = strconv.Atoi(maxStr)
			}
			i.lastResult = rand.Intn(max-min+1) + min
		case 26: // len
			varName := params["var"]
			if val, ok := i.variables[varName].(string); ok {
				i.lastResult = utf8.RuneCountInString(val)
			} else {
				fmt.Println("Ошибка: переменная не является строкой")
			}
		case 27: // substr
			varName := params["var"]
			startStr := params["start"]
			lengthStr := params["length"]
			start, _ := strconv.Atoi(startStr)
			length, _ := strconv.Atoi(lengthStr)
			if val, ok := i.variables[varName].(string); ok {
				if start >= 0 && start < len(val) && length >= 0 {
					if start+length > len(val) {
						length = len(val) - start
					}
					i.lastResult = val[start : start+length]
				} else {
					fmt.Println("Ошибка: неверные индексы")
				}
			} else {
				fmt.Println("Ошибка: переменная не является строкой")
			}
		case 28: // find
			strVar := params["str"]
			subVar := params["sub"]
			if str, ok := i.variables[strVar].(string); ok {
				if sub, ok := i.variables[subVar].(string); ok {
					i.lastResult = strings.Index(str, sub)
				} else {
					fmt.Println("Ошибка: подстрока не является строкой")
				}
			} else {
				fmt.Println("Ошибка: строка не является строкой")
			}
		case 29: // replace
			strVar := params["str"]
			oldVar := params["old"]
			newVar := params["new"]
			if str, ok := i.variables[strVar].(string); ok {
				if old, ok := i.variables[oldVar].(string); ok {
					if new, ok := i.variables[newVar].(string); ok {
						i.lastResult = strings.Replace(str, old, new, -1)
					} else {
						fmt.Println("Ошибка: новое значение не является строкой")
					}
				} else {
					fmt.Println("Ошибка: старое значение не является строкой")
				}
			} else {
				fmt.Println("Ошибка: строка не является строкой")
			}
		case 30: // split
			strVar := params["str"]
			sepVar := params["sep"]
			if str, ok := i.variables[strVar].(string); ok {
				if sep, ok := i.variables[sepVar].(string); ok {
					i.lastResult = strings.Split(str, sep)
				} else {
					fmt.Println("Ошибка: разделитель не является строкой")
				}
			} else {
				fmt.Println("Ошибка: строка не является строкой")
			}
		case 31: // join
			sliceVar := params["slice"]
			sepVar := params["sep"]
			if slice, ok := i.variables[sliceVar].([]string); ok {
				if sep, ok := i.variables[sepVar].(string); ok {
					i.lastResult = strings.Join(slice, sep)
				} else {
					fmt.Println("Ошибка: разделитель не является строкой")
				}
			} else {
				fmt.Println("Ошибка: переменная не является срезом строк")
			}
		case 32: // lower
			varName := params["var"]
			if val, ok := i.variables[varName].(string); ok {
				i.lastResult = strings.ToLower(val)
			} else {
				fmt.Println("Ошибка: переменная не является строкой")
			}
		case 33: // for
			varName := params["var"]
			startStr := params["start"]
			endStr := params["end"]
			start, _ := strconv.Atoi(startStr)
			end, _ := strconv.Atoi(endStr)
			for j := start; j <= end; j++ {
				i.variables[varName] = j
				for _, stmt := range i.ifBlock {
					i.ExecuteStatement(stmt)
				}
			}
			i.ifBlock = []string{}
		case 34: // while
			condVar := params["var"]
			condValueStr := params["value"]
			condValue, _ := strconv.Atoi(condValueStr)
			for {
				if val, ok := i.variables[condVar].(int); ok && val == condValue {
					for _, stmt := range i.ifBlock {
						i.ExecuteStatement(stmt)
					}
				} else {
					break
				}
			}
			i.ifBlock = []string{}
		case 35: // do
			i.inIfBlock = true
		case 36: // while_do
			condVar := params["var"]
			condValueStr := params["value"]
			condValue, _ := strconv.Atoi(condValueStr)
			for {
				for _, stmt := range i.ifBlock {
					i.ExecuteStatement(stmt)
				}
				if val, ok := i.variables[condVar].(int); ok && val != condValue {
					break
				}
			}
			i.ifBlock = []string{}
		case 37: // switch
			varName := params["var"]
			i.variables["switch_var"] = i.variables[varName]
			i.inIfBlock = true
		case 38: // case
			if i.inIfBlock {
				valueStr := params["value"]
				value, _ := strconv.Atoi(valueStr)
				if val, ok := i.variables["switch_var"].(int); ok && val == value {
					for _, stmt := range i.ifBlock {
						i.ExecuteStatement(stmt)
					}
				}
			}
			i.ifBlock = []string{}
		case 39: // default
			if i.inIfBlock {
				for _, stmt := range i.ifBlock {
					i.ExecuteStatement(stmt)
				}
			}
			i.ifBlock = []string{}
		case 40: // file.read
			fileName := params["file"]
			content, err := os.ReadFile(fileName)
			if err != nil {
				fmt.Println("Ошибка чтения файла:", err)
			} else {
				i.lastResult = string(content)
			}
		case 41: // file.write
			fileName := params["file"]
			if content, ok := i.lastResult.(string); ok {
				err := os.WriteFile(fileName, []byte(content), 0644)
				if err != nil {
					fmt.Println("Ошибка записи файла:", err)
				}
			}
		case 42: // print_formatted
			varName := params["var"]
			if val, ok := i.variables[varName]; ok {
				fmt.Printf("%v\n", val)
			} else {
				fmt.Printf("%s\n", varName)
			}
		case 43: // input
			varName := params["var"]
			fmt.Printf("Введите значение для %s: ", varName)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			i.variables[varName] = strings.TrimSpace(input)
		case 44: // array_create
			sizeStr := params["size"]
			size, _ := strconv.Atoi(sizeStr)
			i.lastResult = make([]interface{}, size)
		case 45: // array_set
			arrayVar := params["array"]
			indexStr := params["index"]
			valueVar := params["value"]
			index, _ := strconv.Atoi(indexStr)
			if array, ok := i.variables[arrayVar].([]interface{}); ok {
				if index >= 0 && index < len(array) {
					if val, ok := i.variables[valueVar]; ok {
						array[index] = val
					}
				}
			}
		case 46: // array_get
			arrayVar := params["array"]
			indexStr := params["index"]
			index, _ := strconv.Atoi(indexStr)
			if array, ok := i.variables[arrayVar].([]interface{}); ok {
				if index >= 0 && index < len(array) {
					i.lastResult = array[index]
				} else {
					fmt.Println("Ошибка: индекс вне диапазона")
				}
			}
		case 47: // list_create
			i.lastResult = []interface{}{}
		case 48: // list_append
			listVar := params["list"]
			valueVar := params["value"]
			if list, ok := i.variables[listVar].([]interface{}); ok {
				if val, ok := i.variables[valueVar]; ok {
					i.variables[listVar] = append(list, val)
				}
			}
		case 49: // list_get
			listVar := params["list"]
			indexStr := params["index"]
			index, _ := strconv.Atoi(indexStr)
			if list, ok := i.variables[listVar].([]interface{}); ok {
				if index >= 0 && index < len(list) {
					i.lastResult = list[index]
				} else {
					fmt.Println("Ошибка: индекс вне диапазона")
				}
			}
		case 50: // dict_create
			i.lastResult = make(map[string]interface{})
		case 51: // dict_set
			dictVar := params["dict"]
			keyVar := params["key"]
			valueVar := params["value"]
			if dict, ok := i.variables[dictVar].(map[string]interface{}); ok {
				if key, ok := i.variables[keyVar].(string); ok {
					if val, ok := i.variables[valueVar]; ok {
						dict[key] = val
					}
				}
			}
		case 52: // dict_get
			dictVar := params["dict"]
			keyVar := params["key"]
			if dict, ok := i.variables[dictVar].(map[string]interface{}); ok {
				if key, ok := i.variables[keyVar].(string); ok {
					if val, exists := dict[key]; exists {
						i.lastResult = val
					} else {
						fmt.Println("Ошибка: ключ не найден")
					}
				}
			}
		case 53: // time
			i.lastResult = time.Now().Format("15:04:05")
		case 54: // date
			i.lastResult = time.Now().Format("2006-01-02")
		case 55: // env
			varName := params["var"]
			i.lastResult = os.Getenv(varName)
		case 56: // def
			funcName := params["name"]
			i.functions[funcName] = []string{}
		case 57: // function_call
			funcName := params["func"]
			if cmds, ok := i.functions[funcName]; ok {
				for _, cmd := range cmds {
					i.ExecuteStatement(cmd)
				}
			}
		}
	}
}

func parseExpression(expr string, variables map[string]interface{}) interface{} {
	expr = strings.ReplaceAll(expr, " ", "")
	parts := strings.Split(expr, "+")
	if len(parts) == 2 {
		left := getValue(parts[0], variables)
		right := getValue(parts[1], variables)
		if leftFloat, ok := left.(float64); ok {
			if rightFloat, ok := right.(float64); ok {
				return leftFloat + rightFloat
			}
		}
		return left.(int) + right.(int)
	}
	parts = strings.Split(expr, "-")
	if len(parts) == 2 {
		left := getValue(parts[0], variables)
		right := getValue(parts[1], variables)
		if leftFloat, ok := left.(float64); ok {
			if rightFloat, ok := right.(float64); ok {
				return leftFloat - rightFloat
			}
		}
		return left.(int) - right.(int)
	}
	parts = strings.Split(expr, "*")
	if len(parts) == 2 {
		left := getValue(parts[0], variables)
		right := getValue(parts[1], variables)
		if leftFloat, ok := left.(float64); ok {
			if rightFloat, ok := right.(float64); ok {
				return leftFloat * rightFloat
			}
		}
		return left.(int) * right.(int)
	}
	parts = strings.Split(expr, "/")
	if len(parts) == 2 {
		left := getValue(parts[0], variables)
		right := getValue(parts[1], variables)
		if leftFloat, ok := left.(float64); ok {
			if rightFloat, ok := right.(float64); ok {
				if rightFloat == 0 {
					fmt.Println("Ошибка: деление на ноль")
					return 0
				}
				return leftFloat / rightFloat
			}
		}
		if right.(int) == 0 {
			fmt.Println("Ошибка: деление на ноль")
			return 0
		}
		return left.(int) / right.(int)
	}
	return getValue(expr, variables)
}

func getValue(part string, variables map[string]interface{}) interface{} {
	if val, ok := variables[part]; ok {
		return val
	}
	if num, err := strconv.Atoi(part); err == nil {
		return num
	}
	if num, err := strconv.ParseFloat(part, 64); err == nil {
		return num
	}
	return part
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
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lineLower := strings.ToLower(line)

		switch {
		case strings.HasPrefix(lineLower, "function ("):
			currentFunction = strings.Trim(strings.TrimPrefix(line, "Function ("), ")")
			i.functions[currentFunction] = []string{}
		case lineLower == "memory start (":
			continue
		case lineLower == ")":
			currentFunction = ""
		case currentFunction != "":
			i.functions[currentFunction] = append(i.functions[currentFunction], line)
		default:
			i.ExecuteStatement(line)
		}
	}
}
