package io

import (
	"bufio"
	"fmt"
	"os"
)

type IOModule struct{}

func NewIOModule() *IOModule {
	return &IOModule{}
}

func (i *IOModule) HandlePrint(cmd Command, params map[string]string, variables map[string]interface{}) {
	varName := params["var"]
	if val, ok := variables[varName]; ok {
		fmt.Println(val)
	} else {
		fmt.Println(varName)
	}
}

func (i *IOModule) HandleSolveInput(cmd Command, params map[string]string, variables map[string]interface{}) {
	varName := params["var"]
	fmt.Printf("Введите число для %s: ", varName)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	num, _ := strconv.Atoi(strings.TrimSpace(input))
	variables[varName] = num
}

func (i *IOModule) HandleTextInput(cmd Command, params map[string]string, variables map[string]interface{}) {
	varName := params["var"]
	fmt.Printf("Введите текст для %s: ", varName)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	variables[varName] = strings.TrimSpace(input)
}

func (i *IOModule) HandlePrintFormatted(cmd Command, params map[string]string, variables map[string]interface{}) {
	varName := params["var"]
	if val, ok := variables[varName]; ok {
		fmt.Printf("%v\n", val)
	} else {
		fmt.Printf("%s\n", varName)
	}
}

func (i *IOModule) HandleInput(cmd Command, params map[string]string, variables map[string]interface{}) {
	varName := params["var"]
	fmt.Printf("Введите значение для %s: ", varName)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	variables[varName] = strings.TrimSpace(input)
}

func (i *IOModule) HandleFileRead(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	fileName := params["file"]
	content, err := os.ReadFile(fileName)
	if err != nil {
		*lastResult = ""
		fmt.Println("Ошибка чтения файла:", err)
	} else {
		*lastResult = string(content)
	}
}

func (i *IOModule) HandleFileWrite(cmd Command, params map[string]string, variables map[string]interface{}, lastResult interface{}) {
	fileName := params["file"]
	if content, ok := lastResult.(string); ok {
		err := os.WriteFile(fileName, []byte(content), 0644)
		if err != nil {
			fmt.Println("Ошибка записи файла:", err)
		}
	}
}