package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Interpreter struct {
	memory      *Memory
	fileManager *FileManager
	lastResult  interface{}
	currentFunc string
	inBlock     bool
	blockLines  []string
	lastIfTrue  bool
	inElse      bool
	pendingIf   string
	reader      *bufio.Reader
	inRepeat    bool
	repeatCount int
	repeatLines []string
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		memory:      NewMemory(),
		fileManager: NewFileManager(),
		reader:      bufio.NewReader(os.Stdin),
	}
}

func (i *Interpreter) ExecuteProgram(program string) {
	lines := strings.Split(program, "\n")

	for _, line := range lines {
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lineLower := strings.ToLower(line)

		// function(name) or function(name) {
		if strings.HasPrefix(lineLower, "function(") || strings.HasPrefix(lineLower, "function (") {
			name := extractBetween(line, "function(", ")")
			if name == "" {
				name = extractBetween(line, "function (", ")")
			}
			i.currentFunc = strings.TrimSpace(name)
			if strings.Contains(line, "{") {
				i.inBlock = true
				i.blockLines = []string{}
			}
			continue
		}

		// memory.start{ or memory.start {
		if strings.HasPrefix(lineLower, "memory.start{") || strings.HasPrefix(lineLower, "memory.start {") ||
			strings.HasPrefix(lineLower, "memory start{") || strings.HasPrefix(lineLower, "memory start {") {
			i.inBlock = true
			i.blockLines = []string{}
			continue
		}

		// repeat(n){ or repeat(n) {
		if strings.HasPrefix(lineLower, "repeat(") || strings.HasPrefix(lineLower, "repeat (") {
			nStr := extractBetween(line, "repeat(", "){")
			if nStr == "" {
				nStr = extractBetween(line, "repeat(", ") {")
			}
			if nStr == "" {
				nStr = extractBetween(line, "repeat (", "){")
			}
			if nStr == "" {
				nStr = extractBetween(line, "repeat (", ") {")
			}
			nStr = strings.TrimSpace(nStr)
			n := 0
			if num, err := strconv.Atoi(nStr); err == nil {
				n = num
			} else if val, ok := i.memory.Get(nStr); ok {
				if num, ok := val.(int); ok {
					n = num
				}
			}
			i.repeatCount = n
			i.repeatLines = []string{}
			i.inRepeat = true
			continue
		}

		// } end of block or repeat
		if line == "}" && i.inBlock {
			if i.currentFunc != "" {
				i.memory.StoreFunc(i.currentFunc, i.blockLines)
			}
			i.inBlock = false
			i.currentFunc = ""
			i.blockLines = []string{}
			continue
		}

		// } end of repeat
		if line == "}" && i.inRepeat {
			for j := 0; j < i.repeatCount; j++ {
				for _, cmd := range i.repeatLines {
					i.ExecuteStatement(cmd)
				}
			}
			i.inRepeat = false
			i.repeatCount = 0
			i.repeatLines = []string{}
			continue
		}

		// memory.load(name1, name2)
		if strings.HasPrefix(lineLower, "memory.load(") || strings.HasPrefix(lineLower, "memory.load (") ||
			strings.HasPrefix(lineLower, "memory load(") || strings.HasPrefix(lineLower, "memory load (") {
			names := extractBetween(line, "memory.load(", ")")
			if names == "" {
				names = extractBetween(line, "memory.load (", ")")
			}
			if names == "" {
				names = extractBetween(line, "memory load(", ")")
			}
			if names == "" {
				names = extractBetween(line, "memory load (", ")")
			}
			funcNames := strings.Split(names, ",")
			for _, fn := range funcNames {
				fn = strings.TrimSpace(fn)
				if cmds, ok := i.memory.GetFunc(fn); ok {
					for _, cmd := range cmds {
						i.ExecuteStatement(cmd)
					}
				}
			}
			continue
		}

		// memory.out
		if lineLower == "memory.out" || lineLower == "memory out" {
			i.memory.PrintAll()
			continue
		}

		if i.inBlock {
			i.blockLines = append(i.blockLines, line)
		} else if i.inRepeat {
			i.repeatLines = append(i.repeatLines, line)
		} else {
			i.ExecuteStatement(line)
		}
	}
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

	// else:
	if lineLower == "else:" || lineLower == "else :" {
		if !i.lastIfTrue {
			i.inElse = true
		}
		return
	}

	if i.inElse {
		i.inElse = false
		i.ExecuteStatement(line)
		return
	}

	// then (continuation of if)
	if strings.HasPrefix(lineLower, "then ") && i.pendingIf != "" {
		action := strings.TrimSpace(strings.TrimPrefix(line, "then "))
		i.executeIf(i.pendingIf, action)
		i.pendingIf = ""
		return
	}

	// Variable assignment: name = value
	if strings.Contains(line, " = ") && !strings.HasPrefix(lineLower, "if ") &&
		!strings.HasPrefix(lineLower, "solve.out") && !strings.HasPrefix(lineLower, "text.out") &&
		!strings.HasPrefix(lineLower, "file.text.copy") && !strings.HasPrefix(lineLower, "solve.input") &&
		!strings.HasPrefix(lineLower, "text.input") &&
		!strings.HasPrefix(lineLower, "upper.out") && !strings.HasPrefix(lineLower, "lower.out") &&
		!strings.HasPrefix(lineLower, "len.out") && !strings.HasPrefix(lineLower, "reverse.out") &&
		!strings.HasPrefix(lineLower, "replace.out") && !strings.HasPrefix(lineLower, "substr.out") &&
		!strings.HasPrefix(lineLower, "find.out") && !strings.HasPrefix(lineLower, "abs.out") &&
		!strings.HasPrefix(lineLower, "random.out") && !strings.HasPrefix(lineLower, "type.out") {
		parts := strings.SplitN(line, " = ", 2)
		if len(parts) == 2 {
			varName := strings.TrimSpace(parts[0])
			valueStr := strings.TrimSpace(parts[1])
			if num, err := strconv.Atoi(valueStr); err == nil {
				i.memory.Set(varName, num)
			} else if num, err := strconv.ParseFloat(valueStr, 64); err == nil {
				i.memory.Set(varName, num)
			} else {
				i.memory.Set(varName, valueStr)
			}
			return
		}
	}

	// print(var)
	if strings.HasPrefix(lineLower, "print(") || strings.HasPrefix(lineLower, "print (") {
		varName := extractBetween(line, "print(", ")")
		if varName == "" {
			varName = extractBetween(line, "print (", ")")
		}
		varName = strings.TrimSpace(varName)
		if val, ok := i.memory.Get(varName); ok {
			fmt.Println(val)
		} else {
			fmt.Println(varName)
		}
		return
	}

	// pass()
	if lineLower == "pass()" || lineLower == "pass ()" {
		return
	}

	// jump(name)
	if strings.HasPrefix(lineLower, "jump(") || strings.HasPrefix(lineLower, "jump (") {
		funcName := extractBetween(line, "jump(", ")")
		if funcName == "" {
			funcName = extractBetween(line, "jump (", ")")
		}
		i.executeFunction(strings.TrimSpace(funcName))
		return
	}

	// memory.start(name) - call function
	if strings.HasPrefix(lineLower, "memory.start(") || strings.HasPrefix(lineLower, "memory.start (") ||
		strings.HasPrefix(lineLower, "memory start(") || strings.HasPrefix(lineLower, "memory start (") {
		funcName := extractBetween(line, "memory.start(", ")")
		if funcName == "" {
			funcName = extractBetween(line, "memory.start (", ")")
		}
		if funcName == "" {
			funcName = extractBetween(line, "memory start(", ")")
		}
		if funcName == "" {
			funcName = extractBetween(line, "memory start (", ")")
		}
		i.executeFunction(strings.TrimSpace(funcName))
		return
	}

	// debug() — print all variables and functions
	if lineLower == "debug()" || lineLower == "debug ()" {
		i.memory.PrintDebug()
		return
	}

	// decode() — decode memory.out values
	if lineLower == "decode()" || lineLower == "decode ()" {
		i.memory.DecodePrint()
		return
	}

	// assert(var, value) — assertion for testing
	if strings.HasPrefix(lineLower, "assert(") || strings.HasPrefix(lineLower, "assert (") {
		args := extractBetween(line, "assert(", ")")
		if args == "" {
			args = extractBetween(line, "assert (", ")")
		}
		parts := strings.SplitN(args, ",", 2)
		if len(parts) == 2 {
			varName := strings.TrimSpace(parts[0])
			valueStr := strings.TrimSpace(parts[1])
			if val, ok := i.memory.Get(varName); ok {
				if intVal, ok := val.(int); ok {
					if targetVal, err := strconv.Atoi(valueStr); err == nil {
						if intVal != targetVal {
							fmt.Printf("ASSERT FAIL: %s = %d, expected %d\n", varName, intVal, targetVal)
						} else {
							fmt.Printf("ASSERT OK: %s = %d\n", varName, intVal)
						}
					}
				} else if strVal, ok := val.(string); ok {
					if strVal != valueStr {
						fmt.Printf("ASSERT FAIL: %s = %s, expected %s\n", varName, strVal, valueStr)
					} else {
						fmt.Printf("ASSERT OK: %s = %s\n", varName, strVal)
					}
				}
			}
		}
		return
	}

	// len(var) — string length
	if strings.HasPrefix(lineLower, "len(") || strings.HasPrefix(lineLower, "len (") {
		varName := extractBetween(line, "len(", ")")
		if varName == "" {
			varName = extractBetween(line, "len (", ")")
		}
		varName = strings.TrimSpace(varName)
		if val, ok := i.memory.Get(varName); ok {
			if str, ok := val.(string); ok {
				i.lastResult = len([]rune(str))
			}
		}
		return
	}

	// upper(var) — uppercase
	if strings.HasPrefix(lineLower, "upper(") || strings.HasPrefix(lineLower, "upper (") {
		varName := extractBetween(line, "upper(", ")")
		if varName == "" {
			varName = extractBetween(line, "upper (", ")")
		}
		varName = strings.TrimSpace(varName)
		if val, ok := i.memory.Get(varName); ok {
			if str, ok := val.(string); ok {
				i.lastResult = strings.ToUpper(str)
			}
		}
		return
	}

	// lower(var) — lowercase
	if strings.HasPrefix(lineLower, "lower(") || strings.HasPrefix(lineLower, "lower (") {
		varName := extractBetween(line, "lower(", ")")
		if varName == "" {
			varName = extractBetween(line, "lower (", ")")
		}
		varName = strings.TrimSpace(varName)
		if val, ok := i.memory.Get(varName); ok {
			if str, ok := val.(string); ok {
				i.lastResult = strings.ToLower(str)
			}
		}
		return
	}

	// type(var) — get type of variable
	if strings.HasPrefix(lineLower, "type(") || strings.HasPrefix(lineLower, "type (") {
		varName := extractBetween(line, "type(", ")")
		if varName == "" {
			varName = extractBetween(line, "type (", ")")
		}
		varName = strings.TrimSpace(varName)
		if val, ok := i.memory.Get(varName); ok {
			switch val.(type) {
			case int:
				i.lastResult = "int"
			case float64:
				i.lastResult = "float"
			case string:
				i.lastResult = "string"
			case bool:
				i.lastResult = "bool"
			default:
				i.lastResult = "unknown"
			}
		}
		return
	}

	// abs(var) — absolute value
	if strings.HasPrefix(lineLower, "abs(") || strings.HasPrefix(lineLower, "abs (") {
		varName := extractBetween(line, "abs(", ")")
		if varName == "" {
			varName = extractBetween(line, "abs (", ")")
		}
		varName = strings.TrimSpace(varName)
		if val, ok := i.memory.Get(varName); ok {
			if num, ok := toFloat(val); ok {
				i.lastResult = math.Abs(num)
			}
		}
		return
	}

	// random() — random number 0-99
	if lineLower == "random()" || lineLower == "random ()" {
		i.lastResult = rand.Intn(100)
		return
	}

	// replace(var, old, new) — string replace
	if strings.HasPrefix(lineLower, "replace(") || strings.HasPrefix(lineLower, "replace (") {
		args := extractBetween(line, "replace(", ")")
		if args == "" {
			args = extractBetween(line, "replace (", ")")
		}
		parts := strings.SplitN(args, ",", 3)
		if len(parts) == 3 {
			varName := strings.TrimSpace(parts[0])
			oldStr := strings.TrimSpace(parts[1])
			newStr := strings.TrimSpace(parts[2])
			if val, ok := i.memory.Get(varName); ok {
				if str, ok := val.(string); ok {
					i.lastResult = strings.ReplaceAll(str, oldStr, newStr)
				}
			}
		}
		return
	}

	// substr(var, start, len) — substring
	if strings.HasPrefix(lineLower, "substr(") || strings.HasPrefix(lineLower, "substr (") {
		args := extractBetween(line, "substr(", ")")
		if args == "" {
			args = extractBetween(line, "substr (", ")")
		}
		parts := strings.SplitN(args, ",", 3)
		if len(parts) == 3 {
			varName := strings.TrimSpace(parts[0])
			startStr := strings.TrimSpace(parts[1])
			lenStr := strings.TrimSpace(parts[2])
			if val, ok := i.memory.Get(varName); ok {
				if str, ok := val.(string); ok {
					start, _ := strconv.Atoi(startStr)
					length, _ := strconv.Atoi(lenStr)
					runes := []rune(str)
					if start >= 0 && start < len(runes) && start+length <= len(runes) {
						i.lastResult = string(runes[start : start+length])
					}
				}
			}
		}
		return
	}

	// find(var, sub) — find substring position
	if strings.HasPrefix(lineLower, "find(") || strings.HasPrefix(lineLower, "find (") {
		args := extractBetween(line, "find(", ")")
		if args == "" {
			args = extractBetween(line, "find (", ")")
		}
		parts := strings.SplitN(args, ",", 2)
		if len(parts) == 2 {
			varName := strings.TrimSpace(parts[0])
			subStr := strings.TrimSpace(parts[1])
			if val, ok := i.memory.Get(varName); ok {
				if str, ok := val.(string); ok {
					i.lastResult = strings.Index(str, subStr)
				}
			}
		}
		return
	}

	// reverse(var) — reverse string
	if strings.HasPrefix(lineLower, "reverse(") || strings.HasPrefix(lineLower, "reverse (") {
		varName := extractBetween(line, "reverse(", ")")
		if varName == "" {
			varName = extractBetween(line, "reverse (", ")")
		}
		varName = strings.TrimSpace(varName)
		if val, ok := i.memory.Get(varName); ok {
			if str, ok := val.(string); ok {
				runes := []rune(str)
				for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
					runes[i], runes[j] = runes[j], runes[i]
				}
				i.lastResult = string(runes)
			}
		}
		return
	}

	// input() — general input (auto-detect)
	if lineLower == "input()" || lineLower == "input ()" {
		fmt.Print("Input: ")
		input, _ := i.reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if num, err := strconv.Atoi(input); err == nil {
			i.lastResult = num
		} else {
			i.lastResult = input
		}
		return
	}

	// clear() — clear all variables
	if lineLower == "clear()" || lineLower == "clear ()" {
		i.memory.ClearVars()
		return
	}

	// if
	if strings.HasPrefix(lineLower, "if ") {
		i.handleIf(line)
		return
	}

	// solve.input() = var
	if strings.HasPrefix(lineLower, "solve.input() = ") {
		varName := strings.TrimSpace(strings.TrimPrefix(lineLower, "solve.input() = "))
		fmt.Printf("Enter number for %s: ", varName)
		input, _ := i.reader.ReadString('\n')
		num, _ := strconv.Atoi(strings.TrimSpace(input))
		i.memory.Set(varName, num)
		return
	}

	// solve.out = var
	if strings.HasPrefix(lineLower, "solve.out = ") {
		varName := strings.TrimSpace(strings.TrimPrefix(lineLower, "solve.out = "))
		i.memory.Set(varName, i.lastResult)
		return
	}

	// upper.out = var
	if strings.HasPrefix(lineLower, "upper.out = ") {
		varName := strings.TrimSpace(line[len("upper.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// lower.out = var
	if strings.HasPrefix(lineLower, "lower.out = ") {
		varName := strings.TrimSpace(line[len("lower.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// len.out = var
	if strings.HasPrefix(lineLower, "len.out = ") {
		varName := strings.TrimSpace(line[len("len.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// reverse.out = var
	if strings.HasPrefix(lineLower, "reverse.out = ") {
		varName := strings.TrimSpace(line[len("reverse.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// replace.out = var
	if strings.HasPrefix(lineLower, "replace.out = ") {
		varName := strings.TrimSpace(line[len("replace.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// substr.out = var
	if strings.HasPrefix(lineLower, "substr.out = ") {
		varName := strings.TrimSpace(line[len("substr.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// find.out = var
	if strings.HasPrefix(lineLower, "find.out = ") {
		varName := strings.TrimSpace(line[len("find.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// abs.out = var
	if strings.HasPrefix(lineLower, "abs.out = ") {
		varName := strings.TrimSpace(line[len("abs.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// random.out = var
	if strings.HasPrefix(lineLower, "random.out = ") {
		varName := strings.TrimSpace(line[len("random.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// type.out = var
	if strings.HasPrefix(lineLower, "type.out = ") {
		varName := strings.TrimSpace(line[len("type.out = "):])
		i.memory.Set(varName, i.lastResult)
		return
	}

	// solve(expr)
	if strings.HasPrefix(lineLower, "solve(") || strings.HasPrefix(lineLower, "solve (") {
		expr := extractBetween(line, "solve(", ")")
		if expr == "" {
			expr = extractBetween(line, "solve (", ")")
		}
		i.lastResult = i.evaluateMath(expr)
		return
	}

	// text.input() = var
	if strings.HasPrefix(lineLower, "text.input() = ") {
		varName := strings.TrimSpace(strings.TrimPrefix(lineLower, "text.input() = "))
		fmt.Printf("Enter text for %s: ", varName)
		input, _ := i.reader.ReadString('\n')
		i.memory.Set(varName, strings.TrimSpace(input))
		return
	}

	// text.out = var
	if strings.HasPrefix(lineLower, "text.out = ") {
		varName := strings.TrimSpace(strings.TrimPrefix(lineLower, "text.out = "))
		i.memory.Set(varName, i.lastResult)
		return
	}

	// text(expr)
	if strings.HasPrefix(lineLower, "text(") || strings.HasPrefix(lineLower, "text (") {
		expr := extractBetween(line, "text(", ")")
		if expr == "" {
			expr = extractBetween(line, "text (", ")")
		}
		i.lastResult = i.evaluateText(expr)
		return
	}

	// file operations
	if strings.HasPrefix(lineLower, "file.") {
		i.handleFileOp(line)
		return
	}
}

func (i *Interpreter) handleFileOp(line string) {
	lineLower := strings.ToLower(line)

	if strings.HasPrefix(lineLower, "file.open(") || strings.HasPrefix(lineLower, "file.open (") {
		filename := extractBetween(line, "file.open(", ")")
		if filename == "" {
			filename = extractBetween(line, "file.open (", ")")
		}
		i.fileManager.Open(strings.TrimSpace(filename))
		return
	}

	if lineLower == "file.close" {
		i.fileManager.Close()
		return
	}

	if strings.HasPrefix(lineLower, "file.write(") || strings.HasPrefix(lineLower, "file.write (") {
		text := extractBetween(line, "file.write(", ")")
		if text == "" {
			text = extractBetween(line, "file.write (", ")")
		}
		i.fileManager.Write(strings.TrimSpace(text))
		return
	}

	if lineLower == "file.text.read" {
		i.fileManager.Read()
		return
	}

	if strings.HasPrefix(lineLower, "file.text.copy = ") {
		varName := strings.TrimSpace(strings.TrimPrefix(lineLower, "file.text.copy = "))
		i.memory.Set(varName, i.fileManager.GetContent())
		return
	}

	if strings.HasPrefix(lineLower, "file.create(") || strings.HasPrefix(lineLower, "file.create (") {
		filename := extractBetween(line, "file.create(", ")")
		if filename == "" {
			filename = extractBetween(line, "file.create (", ")")
		}
		i.fileManager.Create(strings.TrimSpace(filename))
		return
	}
}

func (i *Interpreter) handleIf(line string) {
	lineLower := strings.ToLower(line)
	i.inElse = false

	if strings.HasPrefix(lineLower, "if file.none = ") {
		action := strings.TrimSpace(strings.TrimPrefix(line, "if file.none = "))
		if !i.fileManager.IsOpen() {
			actions := strings.Split(action, ",")
			for _, act := range actions {
				act = strings.TrimSpace(act)
				i.ExecuteStatement(act)
			}
			i.lastIfTrue = true
		} else {
			i.lastIfTrue = false
		}
		return
	}

	if strings.HasPrefix(lineLower, "if ") {
		if strings.Contains(lineLower, " then ") {
			parts := strings.SplitN(line, " then ", 2)
			if len(parts) == 2 {
				condition := strings.TrimSpace(parts[0][3:])
				action := strings.TrimSpace(parts[1])
				i.executeIf(condition, action)
				return
			}
		} else {
			condition := strings.TrimSpace(strings.TrimPrefix(line, "if "))
			i.pendingIf = condition
			i.lastIfTrue = false
			return
		}
	}
}

func (i *Interpreter) executeIf(condition, action string) {
	if strings.Contains(condition, " = ") {
		condParts := strings.SplitN(condition, " = ", 2)
		if len(condParts) == 2 {
			varName := strings.TrimSpace(condParts[0])
			valueStr := strings.TrimSpace(condParts[1])

			if val, ok := i.memory.Get(varName); ok {
				if intVal, ok := val.(int); ok {
					if targetVal, err := strconv.Atoi(valueStr); err == nil {
						if intVal == targetVal {
							i.lastIfTrue = true
							i.ExecuteStatement(action)
						} else {
							i.lastIfTrue = false
						}
						return
					}
				}
				if strVal, ok := val.(string); ok {
					if strVal == valueStr {
						i.lastIfTrue = true
						i.ExecuteStatement(action)
					} else {
						i.lastIfTrue = false
					}
					return
				}
			}
		}
	}
	i.lastIfTrue = false
}

func (i *Interpreter) evaluateMath(expr string) interface{} {
	expr = strings.ReplaceAll(expr, " ", "")

	if parts := strings.SplitN(expr, "+", 2); len(parts) == 2 {
		left := i.getValue(parts[0])
		right := i.getValue(parts[1])
		if l, ok := toFloat(left); ok {
			if r, ok := toFloat(right); ok {
				return l + r
			}
		}
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l + r
			}
		}
	}

	if parts := strings.SplitN(expr, "-", 2); len(parts) == 2 {
		left := i.getValue(parts[0])
		right := i.getValue(parts[1])
		if l, ok := toFloat(left); ok {
			if r, ok := toFloat(right); ok {
				return l - r
			}
		}
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l - r
			}
		}
	}

	if parts := strings.SplitN(expr, "*", 2); len(parts) == 2 {
		left := i.getValue(parts[0])
		right := i.getValue(parts[1])
		if l, ok := toFloat(left); ok {
			if r, ok := toFloat(right); ok {
				return l * r
			}
		}
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l * r
			}
		}
	}

	if parts := strings.SplitN(expr, "/", 2); len(parts) == 2 {
		left := i.getValue(parts[0])
		right := i.getValue(parts[1])
		if l, ok := toFloat(left); ok {
			if r, ok := toFloat(right); ok {
				if r == 0 {
					fmt.Println("Error: division by zero")
					return 0
				}
				return l / r
			}
		}
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				if r == 0 {
					fmt.Println("Error: division by zero")
					return 0
				}
				return l / r
			}
		}
	}

	return i.getValue(expr)
}

func (i *Interpreter) evaluateText(expr string) string {
	expr = strings.ReplaceAll(expr, " ", "")
	parts := strings.Split(expr, "+")
	var result string
	for _, part := range parts {
		if val, ok := i.memory.Get(part); ok {
			if str, ok := val.(string); ok {
				if result != "" {
					result += " "
				}
				result += str
			}
		} else {
			if result != "" {
				result += " "
			}
			result += part
		}
	}
	return result
}

func (i *Interpreter) getValue(name string) interface{} {
	if val, ok := i.memory.Get(name); ok {
		return val
	}
	if num, err := strconv.Atoi(name); err == nil {
		return num
	}
	if num, err := strconv.ParseFloat(name, 64); err == nil {
		return num
	}
	return name
}

func (i *Interpreter) executeFunction(name string) {
	if cmds, ok := i.memory.GetFunc(name); ok {
		for _, cmd := range cmds {
			i.ExecuteStatement(cmd)
		}
	}
}

func extractBetween(s, start, end string) string {
	sLower := strings.ToLower(s)
	startLower := strings.ToLower(start)
	endLower := strings.ToLower(end)

	startIdx := strings.Index(sLower, startLower)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(start)
	endIdx := strings.Index(sLower[startIdx:], endLower)
	if endIdx == -1 {
		return ""
	}
	return s[startIdx : startIdx+endIdx]
}

func toFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}
