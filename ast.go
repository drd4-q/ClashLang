package main

// NodeType определяет тип узла AST
type NodeType int

const (
	// Выражения
	NodeNumber NodeType = iota
	NodeString
	NodeIdentifier
	NodeBinaryOp
	NodeUnaryOp
	NodeFunctionCall
	NodeList
	NodeDict

	// Операторы
	NodeAssignment
	NodePrint
	NodeIf
	NodeFor
	NodeWhile
	NodeFunction
	NodeReturn
	NodeBreak
	NodeContinue

	// Управление памятью
	NodeMemAlloc
	NodeMemFree
	NodeMemCheck

	// Управление процессами
	NodeProcCreate
	NodeProcRead
	NodeProcKill
	NodeProcList

	// Программа
	NodeProgram
	NodeBlock
)

// Node базовый интерфейс для всех узлов AST
type Node interface {
	Type() NodeType
}

// Program представляет всю программу
type Program struct {
	Statements []Node
}

func (p *Program) Type() NodeType { return NodeProgram }

// NumberNode представляет числовой литерал
type NumberNode struct {
	Value float64
}

func (n *NumberNode) Type() NodeType { return NodeNumber }

// StringNode представляет строковый литерал
type StringNode struct {
	Value string
}

func (s *StringNode) Type() NodeType { return NodeString }

// IdentifierNode представляет идентификатор (имя переменной)
type IdentifierNode struct {
	Name string
}

func (i *IdentifierNode) Type() NodeType { return NodeIdentifier }

// BinaryOpNode представляет бинарную операцию
type BinaryOpNode struct {
	Left     Node
	Operator string
	Right    Node
}

func (b *BinaryOpNode) Type() NodeType { return NodeBinaryOp }

// UnaryOpNode представляет унарную операцию
type UnaryOpNode struct {
	Operator string
	Operand  Node
}

func (u *UnaryOpNode) Type() NodeType { return NodeUnaryOp }

// FunctionCallNode представляет вызов функции
type FunctionCallNode struct {
	Name      string
	Arguments []Node
}

func (f *FunctionCallNode) Type() NodeType { return NodeFunctionCall }

// ListNode представляет список
type ListNode struct {
	Elements []Node
}

func (l *ListNode) Type() NodeType { return NodeList }

// DictNode представляет словарь
type DictNode struct {
	Keys   []Node
	Values []Node
}

func (d *DictNode) Type() NodeType { return NodeDict }

// AssignmentNode представляет присваивание
type AssignmentNode struct {
	Name  string
	Value Node
}

func (a *AssignmentNode) Type() NodeType { return NodeAssignment }

// PrintNode представляет вывод
type PrintNode struct {
	Expression Node
}

func (p *PrintNode) Type() NodeType { return NodePrint }

// IfNode представляет условный оператор
type IfNode struct {
	Condition Node
	ThenBlock []Node
	ElseBlock []Node
}

func (i *IfNode) Type() NodeType { return NodeIf }

// ForNode представляет цикл for
type ForNode struct {
	Variable  string
	Start     Node
	End       Node
	Step      Node
	Body      []Node
}

func (f *ForNode) Type() NodeType { return NodeFor }

// WhileNode представляет цикл while
type WhileNode struct {
	Condition Node
	Body      []Node
}

func (w *WhileNode) Type() NodeType { return NodeWhile }

// FunctionNode представляет определение функции
type FunctionNode struct {
	Name       string
	Parameters []string
	Body       []Node
}

func (f *FunctionNode) Type() NodeType { return NodeFunction }

// ReturnNode представляет возврат из функции
type ReturnNode struct {
	Value Node
}

func (r *ReturnNode) Type() NodeType { return NodeReturn }

// BreakNode представляет break
type BreakNode struct{}

func (b *BreakNode) Type() NodeType { return NodeBreak }

// ContinueNode представляет continue
type ContinueNode struct{}

func (c *ContinueNode) Type() NodeType { return NodeContinue }

// MemAllocNode представляет выделение памяти
type MemAllocNode struct {
	Size Node
	Name string
}

func (m *MemAllocNode) Type() NodeType { return NodeMemAlloc }

// MemFreeNode представляет освобождение памяти
type MemFreeNode struct {
	Name string
}

func (m *MemFreeNode) Type() NodeType { return NodeMemFree }

// MemCheckNode представляет проверку памяти
type MemCheckNode struct{}

func (m *MemCheckNode) Type() NodeType { return NodeMemCheck }

// ProcCreateNode представляет создание процесса
type ProcCreateNode struct {
	Name     string
	CodeFile Node
}

func (p *ProcCreateNode) Type() NodeType { return NodeProcCreate }

// ProcReadNode представляет чтение кода процесса
type ProcReadNode struct {
	ProcID Node
}

func (p *ProcReadNode) Type() NodeType { return NodeProcRead }

// ProcKillNode представляет завершение процесса
type ProcKillNode struct {
	ProcID Node
}

func (p *ProcKillNode) Type() NodeType { return NodeProcKill }

// ProcListNode представляет список процессов
type ProcListNode struct{}

func (p *ProcListNode) Type() NodeType { return NodeProcList }

// BlockNode представляет блок кода
type BlockNode struct {
	Statements []Node
}

func (b *BlockNode) Type() NodeType { return NodeBlock }
