package main

import (
	"fmt"
	"math"
	"os"
)

// NewInterpreter создает новый интерпретатор
type NewInterpreter struct {
	variables      map[string]interface{}
	functions      map[string]*FunctionNode
	memoryManager  *MemoryManager
	processManager *ProcessManager
	returnValue    interface{}
	shouldReturn   bool
	shouldBreak    bool
	shouldContinue bool
}

// NewNewInterpreter создает экземпляр интерпретатора
func NewNewInterpreter() *NewInterpreter {
	return &NewInterpreter{
		variables:      make(map[string]interface{}),
		functions:      make(map[string]*FunctionNode),
		memoryManager:  NewMemoryManager(),
		processManager: NewProcessManager(),
	}
}

// Execute выполняет программу
func (ni *NewInterpreter) Execute(program *Program) error {
	for _, stmt := range program.Statements {
		if err := ni.executeNode(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (ni *NewInterpreter) executeNode(node Node) error {
	if ni.shouldReturn || ni.shouldBreak || ni.shouldContinue {
		return nil
	}

	switch n := node.(type) {
	case *PrintNode:
		return ni.executePrint(n)
	case *AssignmentNode:
		return ni.executeAssignment(n)
	case *IfNode:
		return ni.executeIf(n)
	case *ForNode:
		return ni.executeFor(n)
	case *WhileNode:
		return ni.executeWhile(n)
	case *FunctionNode:
		return ni.executeFunction(n)
	case *FunctionCallNode:
		_, err := ni.executeFunctionCall(n)
		return err
	case *ReturnNode:
		return ni.executeReturn(n)
	case *BreakNode:
		ni.shouldBreak = true
		return nil
	case *ContinueNode:
		ni.shouldContinue = true
		return nil
	case *MemAllocNode:
		return ni.executeMemAlloc(n)
	case *MemFreeNode:
		return ni.executeMemFree(n)
	case *MemCheckNode:
		ni.memoryManager.Check()
		return nil
	case *ProcCreateNode:
		return ni.executeProcCreate(n)
	case *ProcReadNode:
		return ni.executeProcRead(n)
	case *ProcKillNode:
		return ni.executeProcKill(n)
	case *ProcListNode:
		ni.processManager.List()
		return nil
	default:
		return fmt.Errorf("неизвестный тип узла: %T", node)
	}
}

func (ni *NewInterpreter) executePrint(node *PrintNode) error {
	value, err := ni.evaluateExpression(node.Expression)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

func (ni *NewInterpreter) executeAssignment(node *AssignmentNode) error {
	value, err := ni.evaluateExpression(node.Value)
	if err != nil {
		return err
	}
	ni.variables[node.Name] = value
	return nil
}

func (ni *NewInterpreter) executeIf(node *IfNode) error {
	condition, err := ni.evaluateExpression(node.Condition)
	if err != nil {
		return err
	}

	if ni.isTruthy(condition) {
		for _, stmt := range node.ThenBlock {
			if err := ni.executeNode(stmt); err != nil {
				return err
			}
			if ni.shouldReturn || ni.shouldBreak || ni.shouldContinue {
				return nil
			}
		}
	} else if len(node.ElseBlock) > 0 {
		for _, stmt := range node.ElseBlock {
			if err := ni.executeNode(stmt); err != nil {
				return err
			}
			if ni.shouldReturn || ni.shouldBreak || ni.shouldContinue {
				return nil
			}
		}
	}

	return nil
}

func (ni *NewInterpreter) executeFor(node *ForNode) error {
	start, err := ni.evaluateExpression(node.Start)
	if err != nil {
		return err
	}

	end, err := ni.evaluateExpression(node.End)
	if err != nil {
		return err
	}

	step, err := ni.evaluateExpression(node.Step)
	if err != nil {
		return err
	}

	startNum := ni.toFloat(start)
	endNum := ni.toFloat(end)
	stepNum := ni.toFloat(step)

	for i := startNum; i < endNum; i += stepNum {
		ni.variables[node.Variable] = i

		for _, stmt := range node.Body {
			if err := ni.executeNode(stmt); err != nil {
				return err
			}

			if ni.shouldReturn {
				return nil
			}

			if ni.shouldBreak {
				ni.shouldBreak = false
				return nil
			}

			if ni.shouldContinue {
				ni.shouldContinue = false
				break
			}
		}
	}

	return nil
}

func (ni *NewInterpreter) executeWhile(node *WhileNode) error {
	for {
		condition, err := ni.evaluateExpression(node.Condition)
		if err != nil {
			return err
		}

		if !ni.isTruthy(condition) {
			break
		}

		for _, stmt := range node.Body {
			if err := ni.executeNode(stmt); err != nil {
				return err
			}

			if ni.shouldReturn {
				return nil
			}

			if ni.shouldBreak {
				ni.shouldBreak = false
				return nil
			}

			if ni.shouldContinue {
				ni.shouldContinue = false
				break
			}
		}
	}

	return nil
}

func (ni *NewInterpreter) executeFunction(node *FunctionNode) error {
	ni.functions[node.Name] = node
	return nil
}

func (ni *NewInterpreter) executeFunctionCall(node *FunctionCallNode) (interface{}, error) {
	// Встроенные функции
	switch node.Name {
	case "input":
		fmt.Print("Введите значение: ")
		var input string
		fmt.Scanln(&input)
		return input, nil
	case "int":
		if len(node.Arguments) > 0 {
			val, err := ni.evaluateExpression(node.Arguments[0])
			if err != nil {
				return nil, err
			}
			return int(ni.toFloat(val)), nil
		}
	case "float":
		if len(node.Arguments) > 0 {
			val, err := ni.evaluateExpression(node.Arguments[0])
			if err != nil {
				return nil, err
			}
			return ni.toFloat(val), nil
		}
	case "str":
		if len(node.Arguments) > 0 {
			val, err := ni.evaluateExpression(node.Arguments[0])
			if err != nil {
				return nil, err
			}
			return fmt.Sprintf("%v", val), nil
		}
	case "len":
		if len(node.Arguments) > 0 {
			val, err := ni.evaluateExpression(node.Arguments[0])
			if err != nil {
				return nil, err
			}
			if str, ok := val.(string); ok {
				return float64(len(str)), nil
			}
		}
	}

	// Пользовательские функции
	fn, exists := ni.functions[node.Name]
	if !exists {
		return nil, fmt.Errorf("функция '%s' не определена", node.Name)
	}

	if len(node.Arguments) != len(fn.Parameters) {
		return nil, fmt.Errorf("функция '%s' ожидает %d аргументов, получено %d",
			node.Name, len(fn.Parameters), len(node.Arguments))
	}

	// Сохраняем текущие переменные
	oldVars := make(map[string]interface{})
	for k, v := range ni.variables {
		oldVars[k] = v
	}

	// Устанавливаем параметры
	for i, param := range fn.Parameters {
		val, err := ni.evaluateExpression(node.Arguments[i])
		if err != nil {
			return nil, err
		}
		ni.variables[param] = val
	}

	// Выполняем тело функции
	for _, stmt := range fn.Body {
		if err := ni.executeNode(stmt); err != nil {
			return nil, err
		}
		if ni.shouldReturn {
			break
		}
	}

	result := ni.returnValue
	ni.shouldReturn = false
	ni.returnValue = nil

	// Восстанавливаем переменные
	ni.variables = oldVars

	return result, nil
}

func (ni *NewInterpreter) executeReturn(node *ReturnNode) error {
	if node.Value != nil {
		val, err := ni.evaluateExpression(node.Value)
		if err != nil {
			return err
		}
		ni.returnValue = val
	}
	ni.shouldReturn = true
	return nil
}

func (ni *NewInterpreter) executeMemAlloc(node *MemAllocNode) error {
	size, err := ni.evaluateExpression(node.Size)
	if err != nil {
		return err
	}

	sizeInt := int(ni.toFloat(size))
	return ni.memoryManager.Alloc(node.Name, sizeInt)
}

func (ni *NewInterpreter) executeMemFree(node *MemFreeNode) error {
	return ni.memoryManager.Free(node.Name)
}

func (ni *NewInterpreter) executeProcCreate(node *ProcCreateNode) error {
	codeFile, err := ni.evaluateExpression(node.CodeFile)
	if err != nil {
		return err
	}

	codeFilePath := fmt.Sprintf("%v", codeFile)
	code, err := os.ReadFile(codeFilePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла процесса: %v", err)
	}

	procID := ni.processManager.Create(node.Name, string(code))
	ni.variables["last_proc_id"] = float64(procID)
	return nil
}

func (ni *NewInterpreter) executeProcRead(node *ProcReadNode) error {
	procID, err := ni.evaluateExpression(node.ProcID)
	if err != nil {
		return err
	}

	procIDInt := int(ni.toFloat(procID))
	code, err := ni.processManager.Read(procIDInt)
	if err != nil {
		return err
	}

	fmt.Printf("\n=== КОД ПРОЦЕССА #%d ===\n%s\n========================\n", procIDInt, code)
	return nil
}

func (ni *NewInterpreter) executeProcKill(node *ProcKillNode) error {
	procID, err := ni.evaluateExpression(node.ProcID)
	if err != nil {
		return err
	}

	procIDInt := int(ni.toFloat(procID))
	return ni.processManager.Kill(procIDInt)
}

func (ni *NewInterpreter) evaluateExpression(node Node) (interface{}, error) {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value, nil
	case *StringNode:
		return n.Value, nil
	case *IdentifierNode:
		if val, exists := ni.variables[n.Name]; exists {
			return val, nil
		}
		return nil, fmt.Errorf("переменная '%s' не определена", n.Name)
	case *BinaryOpNode:
		return ni.evaluateBinaryOp(n)
	case *UnaryOpNode:
		return ni.evaluateUnaryOp(n)
	case *FunctionCallNode:
		return ni.executeFunctionCall(n)
	case *ListNode:
		elements := make([]interface{}, len(n.Elements))
		for i, elem := range n.Elements {
			val, err := ni.evaluateExpression(elem)
			if err != nil {
				return nil, err
			}
			elements[i] = val
		}
		return elements, nil
	default:
		return nil, fmt.Errorf("неизвестный тип выражения: %T", node)
	}
}

func (ni *NewInterpreter) evaluateBinaryOp(node *BinaryOpNode) (interface{}, error) {
	left, err := ni.evaluateExpression(node.Left)
	if err != nil {
		return nil, err
	}

	right, err := ni.evaluateExpression(node.Right)
	if err != nil {
		return nil, err
	}

	// Строковая конкатенация
	if node.Operator == "+" {
		if leftStr, ok := left.(string); ok {
			if rightStr, ok := right.(string); ok {
				return leftStr + rightStr, nil
			}
		}
	}

	// Числовые операции
	leftNum := ni.toFloat(left)
	rightNum := ni.toFloat(right)

	switch node.Operator {
	case "+":
		return leftNum + rightNum, nil
	case "-":
		return leftNum - rightNum, nil
	case "*":
		return leftNum * rightNum, nil
	case "/":
		if rightNum == 0 {
			return nil, fmt.Errorf("деление на ноль")
		}
		return leftNum / rightNum, nil
	case "%":
		return math.Mod(leftNum, rightNum), nil
	case "**":
		return math.Pow(leftNum, rightNum), nil
	case "==":
		return leftNum == rightNum, nil
	case "!=":
		return leftNum != rightNum, nil
	case "<":
		return leftNum < rightNum, nil
	case ">":
		return leftNum > rightNum, nil
	case "<=":
		return leftNum <= rightNum, nil
	case ">=":
		return leftNum >= rightNum, nil
	default:
		return nil, fmt.Errorf("неизвестный оператор: %s", node.Operator)
	}
}

func (ni *NewInterpreter) evaluateUnaryOp(node *UnaryOpNode) (interface{}, error) {
	operand, err := ni.evaluateExpression(node.Operand)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case "-":
		return -ni.toFloat(operand), nil
	case "not":
		return !ni.isTruthy(operand), nil
	default:
		return nil, fmt.Errorf("неизвестный унарный оператор: %s", node.Operator)
	}
}

func (ni *NewInterpreter) toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		// Попытка преобразовать строку в число
		var f float64
		fmt.Sscanf(v, "%f", &f)
		return f
	default:
		return 0
	}
}

func (ni *NewInterpreter) isTruthy(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	case string:
		return v != ""
	case nil:
		return false
	default:
		return true
	}
}