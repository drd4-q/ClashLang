package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/yourusername/clashlang/math"
	"github.com/yourusername/clashlang/strings"
	"github.com/yourusername/clashlang/control"
	"github.com/yourusername/clashlang/io"
	"github.com/yourusername/clashlang/data"
	"github.com/yourusername/clashlang/system"
	"github.com/yourusername/clashlang/custom"
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
	math       *math.MathModule
	strings    *strings.StringModule
	control    *control.ControlModule
	io         *io.IOModule
	data       *data.DataModule
	system     *system.SystemModule
	custom     *custom.CustomModule
}

func NewInterpreter() *Interpreter {
	interp := &Interpreter{
		commands:  make(map[int]Command),
		variables: make(map[string]interface{}),
		functions: make(map[string][]string),
		ifBlock:   []string{},
		math:      math.NewMathModule(),
		strings:   strings.NewStringModule(),
		control:   control.NewControlModule(),
		io:        io.NewIOModule(),
		data:      data.NewDataModule(),
		system:    system.NewSystemModule(),
		custom:    custom.NewCustomModule(),
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
	case lineLower == "memory out":
		fmt.Println("Global memory:", i.variables)
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
			i.io.HandlePrint(cmd, params, i.variables)
		case 2: // Solve.input
			i.io.HandleSolveInput(cmd, params, i.variables)
		case 3: // Solve
			i.math.HandleSolve(cmd, params, i.variables, &i.lastResult)
		case 4: // Solve.out
			i.math.HandleSolveOut(cmd, params, i.variables, i.lastResult)
		case 5: // Text.input
			i.io.HandleTextInput(cmd, params, i.variables)
		case 6: // Text
			i.strings.HandleText(cmd, params, i.variables, &i.lastResult)
		case 7: // Text.out
			i.strings.HandleTextOut(cmd, params, i.variables, i.lastResult)
		case 8: // If
			i.control.HandleIf(cmd, params, i.variables, &i.inIfBlock, &i.ifBlock)
		case 9: // Jump
			i.custom.HandleJump(cmd, params, i.functions, i.ExecuteStatement)
		case 13: // Text.length
			i.strings.HandleTextLength(cmd, params, i.variables, &i.lastResult)
		case 14: // Text.upper
			i.strings.HandleTextUpper(cmd, params, i.variables, &i.lastResult)
		// Новые функции из модулей
		case 15: // abs
			i.math.HandleAbs(cmd, params, i.variables, &i.lastResult)
		case 16: // sqrt
			i.math.HandleSqrt(cmd, params, i.variables, &i.lastResult)
		case 17: // pow
			i.math.HandlePow(cmd, params, i.variables, &i.lastResult)
		case 18: // round
			i.math.HandleRound(cmd, params, i.variables, &i.lastResult)
		case 19: // sin
			i.math.HandleSin(cmd, params, i.variables, &i.lastResult)
		case 20: // cos
			i.math.HandleCos(cmd, params, i.variables, &i.lastResult)
		case 21: // tan
			i.math.HandleTan(cmd, params, i.variables, &i.lastResult)
		case 22: // log
			i.math.HandleLog(cmd, params, i.variables, &i.lastResult)
		case 23: // log10
			i.math.HandleLog10(cmd, params, i.variables, &i.lastResult)
		case 24: // random
			i.math.HandleRandom(cmd, params, i.variables, &i.lastResult)
		case 25: // randint
			i.math.HandleRandInt(cmd, params, i.variables, &i.lastResult)
		case 26: // len
			i.strings.HandleLen(cmd, params, i.variables, &i.lastResult)
		case 27: // substr
			i.strings.HandleSubstr(cmd, params, i.variables, &i.lastResult)
		case 28: // find
			i.strings.HandleFind(cmd, params, i.variables, &i.lastResult)
		case 29: // replace
			i.strings.HandleReplace(cmd, params, i.variables, &i.lastResult)
		case 30: // split
			i.strings.HandleSplit(cmd, params, i.variables, &i.lastResult)
		case 31: // join
			i.strings.HandleJoin(cmd, params, i.variables, &i.lastResult)
		case 32: // lower
			i.strings.HandleLower(cmd, params, i.variables, &i.lastResult)
		case 33: // for
			i.control.HandleFor(cmd, params, i.variables, i.ExecuteStatement)
		case 34: // while
			i.control.HandleWhile(cmd, params, i.variables, i.ExecuteStatement, &i.inIfBlock)
		case 35: // do
			i.control.HandleDo(cmd, params, i.variables, i.ExecuteStatement, &i.inIfBlock)
		case 36: // while_do
			i.control.HandleWhileDo(cmd, params, i.variables, i.ExecuteStatement, &i.inIfBlock)
		case 37: // switch
			i.control.HandleSwitch(cmd, params, i.variables, i.ExecuteStatement, &i.inIfBlock)
		case 38: // case
			i.control.HandleCase(cmd, params, i.variables, i.ExecuteStatement, &i.inIfBlock)
		case 39: // default
			i.control.HandleDefault(cmd, params, i.variables, i.ExecuteStatement, &i.inIfBlock)
		case 40: // file.read
			i.io.HandleFileRead(cmd, params, i.variables, &i.lastResult)
		case 41: // file.write
			i.io.HandleFileWrite(cmd, params, i.variables, i.lastResult)
		case 42: // print_formatted
			i.io.HandlePrintFormatted(cmd, params, i.variables)
		case 43: // input
			i.io.HandleInput(cmd, params, i.variables)
		case 44: // array_create
			i.data.HandleArrayCreate(cmd, params, i.variables, &i.lastResult)
		case 45: // array_set
			i.data.HandleArraySet(cmd, params, i.variables)
		case 46: // array_get
			i.data.HandleArrayGet(cmd, params, i.variables, &i.lastResult)
		case 47: // list_create
			i.data.HandleListCreate(cmd, params, i.variables, &i.lastResult)
		case 48: // list_append
			i.data.HandleListAppend(cmd, params, i.variables)
		case 49: // list_get
			i.data.HandleListGet(cmd, params, i.variables, &i.lastResult)
		case 50: // dict_create
			i.data.HandleDictCreate(cmd, params, i.variables, &i.lastResult)
		case 51: // dict_set
			i.data.HandleDictSet(cmd, params, i.variables)
		case 52: // dict_get
			i.data.HandleDictGet(cmd, params, i.variables, &i.lastResult)
		case 53: // time
			i.system.HandleTime(cmd, params, i.variables, &i.lastResult)
		case 54: // date
			i.system.HandleDate(cmd, params, i.variables, &i.lastResult)
		case 55: // env
			i.system.HandleEnv(cmd, params, i.variables, &i.lastResult)
		case 56: // def
			i.custom.HandleDef(cmd, params, i.functions)
		case 57: // function_call
			i.custom.HandleFunctionCall(cmd, params, i.functions, i.ExecuteStatement)
		}
	}
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