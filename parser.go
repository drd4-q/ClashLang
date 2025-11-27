package main

import (
	"fmt"
	"strconv"
)

// Parser выполняет синтаксический анализ
type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
}

// NewParser создает новый парсер
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) expectToken(t TokenType) error {
	if p.currentToken.Type != t {
		return fmt.Errorf("ожидался токен %s, получен %s в строке %d",
			TokenTypeString(t), TokenTypeString(p.currentToken.Type), p.currentToken.Line)
	}
	p.nextToken()
	return nil
}

// Parse парсит программу
func (p *Parser) Parse() (*Program, error) {
	program := &Program{Statements: []Node{}}

	for p.currentToken.Type != TokenEOF {
		// Пропускаем пустые строки
		if p.currentToken.Type == TokenNewline {
			p.nextToken()
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program, nil
}

func (p *Parser) parseStatement() (Node, error) {
	switch p.currentToken.Type {
	case TokenPrint:
		return p.parsePrint()
	case TokenIf:
		return p.parseIf()
	case TokenFor:
		return p.parseFor()
	case TokenWhile:
		return p.parseWhile()
	case TokenDef:
		return p.parseFunction()
	case TokenReturn:
		return p.parseReturn()
	case TokenBreak:
		p.nextToken()
		return &BreakNode{}, nil
	case TokenContinue:
		p.nextToken()
		return &ContinueNode{}, nil
	case TokenMemAlloc:
		return p.parseMemAlloc()
	case TokenMemFree:
		return p.parseMemFree()
	case TokenMemCheck:
		p.nextToken()
		return &MemCheckNode{}, nil
	case TokenProcCreate:
		return p.parseProcCreate()
	case TokenProcRead:
		return p.parseProcRead()
	case TokenProcKill:
		return p.parseProcKill()
	case TokenProcList:
		p.nextToken()
		return &ProcListNode{}, nil
	case TokenIdentifier:
		// Может быть присваивание или вызов функции
		if p.peekToken.Type == TokenAssign {
			return p.parseAssignment()
		} else if p.peekToken.Type == TokenLeftParen {
			return p.parseFunctionCall()
		}
		return nil, fmt.Errorf("неожиданный идентификатор '%s' в строке %d",
			p.currentToken.Value, p.currentToken.Line)
	default:
		return nil, fmt.Errorf("неожиданный токен %s в строке %d",
			TokenTypeString(p.currentToken.Type), p.currentToken.Line)
	}
}

func (p *Parser) parsePrint() (Node, error) {
	p.nextToken() // Пропускаем 'print'

	if err := p.expectToken(TokenLeftParen); err != nil {
		return nil, err
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenRightParen); err != nil {
		return nil, err
	}

	return &PrintNode{Expression: expr}, nil
}

func (p *Parser) parseAssignment() (Node, error) {
	name := p.currentToken.Value
	p.nextToken() // Пропускаем имя переменной

	if err := p.expectToken(TokenAssign); err != nil {
		return nil, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &AssignmentNode{Name: name, Value: value}, nil
}

func (p *Parser) parseExpression() (Node, error) {
	return p.parseComparison()
}

func (p *Parser) parseComparison() (Node, error) {
	left, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == TokenEqual || p.currentToken.Type == TokenNotEqual ||
		p.currentToken.Type == TokenLess || p.currentToken.Type == TokenGreater ||
		p.currentToken.Type == TokenLessEq || p.currentToken.Type == TokenGreaterEq {
		op := p.currentToken.Value
		p.nextToken()

		right, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}

		left = &BinaryOpNode{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseAdditive() (Node, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == TokenPlus || p.currentToken.Type == TokenMinus {
		op := p.currentToken.Value
		p.nextToken()

		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}

		left = &BinaryOpNode{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseMultiplicative() (Node, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == TokenMultiply || p.currentToken.Type == TokenDivide ||
		p.currentToken.Type == TokenModulo || p.currentToken.Type == TokenPower {
		op := p.currentToken.Value
		p.nextToken()

		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		left = &BinaryOpNode{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parsePrimary() (Node, error) {
	switch p.currentToken.Type {
	case TokenNumber:
		value, err := strconv.ParseFloat(p.currentToken.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга числа: %v", err)
		}
		p.nextToken()
		return &NumberNode{Value: value}, nil

	case TokenString:
		value := p.currentToken.Value
		p.nextToken()
		return &StringNode{Value: value}, nil

	case TokenIdentifier:
		name := p.currentToken.Value
		p.nextToken()

		// Проверяем, является ли это вызовом функции
		if p.currentToken.Type == TokenLeftParen {
			p.nextToken()
			args := []Node{}

			if p.currentToken.Type != TokenRightParen {
				for {
					arg, err := p.parseExpression()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)

					if p.currentToken.Type != TokenComma {
						break
					}
					p.nextToken()
				}
			}

			if err := p.expectToken(TokenRightParen); err != nil {
				return nil, err
			}

			return &FunctionCallNode{Name: name, Arguments: args}, nil
		}

		return &IdentifierNode{Name: name}, nil

	case TokenLeftParen:
		p.nextToken()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expectToken(TokenRightParen); err != nil {
			return nil, err
		}
		return expr, nil

	case TokenLeftBrace:
		return p.parseList()

	case TokenMinus:
		p.nextToken()
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &UnaryOpNode{Operator: "-", Operand: operand}, nil

	default:
		return nil, fmt.Errorf("неожиданный токен %s в выражении (строка %d)",
			TokenTypeString(p.currentToken.Type), p.currentToken.Line)
	}
}

func (p *Parser) parseList() (Node, error) {
	p.nextToken() // Пропускаем '['

	elements := []Node{}

	if p.currentToken.Type != TokenRightBrace {
		for {
			elem, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			elements = append(elements, elem)

			if p.currentToken.Type != TokenComma {
				break
			}
			p.nextToken()
		}
	}

	if err := p.expectToken(TokenRightBrace); err != nil {
		return nil, err
	}

	return &ListNode{Elements: elements}, nil
}

func (p *Parser) parseIf() (Node, error) {
	p.nextToken() // Пропускаем 'if'

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenColon); err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenNewline); err != nil {
		return nil, err
	}

	thenBlock, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	var elseBlock []Node
	if p.currentToken.Type == TokenElse {
		p.nextToken()
		if err := p.expectToken(TokenColon); err != nil {
			return nil, err
		}
		if err := p.expectToken(TokenNewline); err != nil {
			return nil, err
		}
		elseBlock, err = p.parseBlock()
		if err != nil {
			return nil, err
		}
	}

	return &IfNode{Condition: condition, ThenBlock: thenBlock, ElseBlock: elseBlock}, nil
}

func (p *Parser) parseFor() (Node, error) {
	p.nextToken() // Пропускаем 'for'

	if p.currentToken.Type != TokenIdentifier {
		return nil, fmt.Errorf("ожидалось имя переменной после 'for'")
	}

	variable := p.currentToken.Value
	p.nextToken()

	if err := p.expectToken(TokenIn); err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenRange); err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenLeftParen); err != nil {
		return nil, err
	}

	start, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenComma); err != nil {
		return nil, err
	}

	end, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	var step Node = &NumberNode{Value: 1}
	if p.currentToken.Type == TokenComma {
		p.nextToken()
		step, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	if err := p.expectToken(TokenRightParen); err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenColon); err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenNewline); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ForNode{Variable: variable, Start: start, End: end, Step: step, Body: body}, nil
}

func (p *Parser) parseWhile() (Node, error) {
	p.nextToken() // Пропускаем 'while'

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenColon); err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenNewline); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &WhileNode{Condition: condition, Body: body}, nil
}

func (p *Parser) parseFunction() (Node, error) {
	p.nextToken() // Пропускаем 'def'

	if p.currentToken.Type != TokenIdentifier {
		return nil, fmt.Errorf("ожидалось имя функции после 'def'")
	}

	name := p.currentToken.Value
	p.nextToken()

	if err := p.expectToken(TokenLeftParen); err != nil {
		return nil, err
	}

	parameters := []string{}
	if p.currentToken.Type != TokenRightParen {
		for {
			if p.currentToken.Type != TokenIdentifier {
				return nil, fmt.Errorf("ожидался параметр функции")
			}
			parameters = append(parameters, p.currentToken.Value)
			p.nextToken()

			if p.currentToken.Type != TokenComma {
				break
			}
			p.nextToken()
		}
	}

	if err := p.expectToken(TokenRightParen); err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenColon); err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenNewline); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &FunctionNode{Name: name, Parameters: parameters, Body: body}, nil
}

func (p *Parser) parseReturn() (Node, error) {
	p.nextToken() // Пропускаем 'return'

	var value Node
	var err error

	if p.currentToken.Type != TokenNewline && p.currentToken.Type != TokenEOF {
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	return &ReturnNode{Value: value}, nil
}

func (p *Parser) parseMemAlloc() (Node, error) {
	p.nextToken() // Пропускаем 'mem_alloc'

	if err := p.expectToken(TokenLeftParen); err != nil {
		return nil, err
	}

	if p.currentToken.Type != TokenIdentifier {
		return nil, fmt.Errorf("ожидалось имя блока памяти")
	}

	name := p.currentToken.Value
	p.nextToken()

	if err := p.expectToken(TokenComma); err != nil {
		return nil, err
	}

	size, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenRightParen); err != nil {
		return nil, err
	}

	return &MemAllocNode{Name: name, Size: size}, nil
}

func (p *Parser) parseMemFree() (Node, error) {
	p.nextToken() // Пропускаем 'mem_free'

	if err := p.expectToken(TokenLeftParen); err != nil {
		return nil, err
	}

	if p.currentToken.Type != TokenIdentifier {
		return nil, fmt.Errorf("ожидалось имя блока памяти")
	}

	name := p.currentToken.Value
	p.nextToken()

	if err := p.expectToken(TokenRightParen); err != nil {
		return nil, err
	}

	return &MemFreeNode{Name: name}, nil
}

func (p *Parser) parseProcCreate() (Node, error) {
	p.nextToken() // Пропускаем 'proc_create'

	if err := p.expectToken(TokenLeftParen); err != nil {
		return nil, err
	}

	if p.currentToken.Type != TokenIdentifier {
		return nil, fmt.Errorf("ожидалось имя процесса")
	}

	name := p.currentToken.Value
	p.nextToken()

	if err := p.expectToken(TokenComma); err != nil {
		return nil, err
	}

	codeFile, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenRightParen); err != nil {
		return nil, err
	}

	return &ProcCreateNode{Name: name, CodeFile: codeFile}, nil
}

func (p *Parser) parseProcRead() (Node, error) {
	p.nextToken() // Пропускаем 'proc_read'

	if err := p.expectToken(TokenLeftParen); err != nil {
		return nil, err
	}

	procID, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenRightParen); err != nil {
		return nil, err
	}

	return &ProcReadNode{ProcID: procID}, nil
}

func (p *Parser) parseProcKill() (Node, error) {
	p.nextToken() // Пропускаем 'proc_kill'

	if err := p.expectToken(TokenLeftParen); err != nil {
		return nil, err
	}

	procID, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(TokenRightParen); err != nil {
		return nil, err
	}

	return &ProcKillNode{ProcID: procID}, nil
}

func (p *Parser) parseFunctionCall() (Node, error) {
	name := p.currentToken.Value
	p.nextToken()

	if err := p.expectToken(TokenLeftParen); err != nil {
		return nil, err
	}

	args := []Node{}
	if p.currentToken.Type != TokenRightParen {
		for {
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if p.currentToken.Type != TokenComma {
				break
			}
			p.nextToken()
		}
	}

	if err := p.expectToken(TokenRightParen); err != nil {
		return nil, err
	}

	return &FunctionCallNode{Name: name, Arguments: args}, nil
}

func (p *Parser) parseBlock() ([]Node, error) {
	statements := []Node{}

	// Простая реализация: читаем все операторы до конца или до уменьшения отступа
	// В реальной реализации нужно отслеживать уровни отступов
	for p.currentToken.Type != TokenEOF {
		// Если встретили токен, который не может быть в блоке, выходим
		if p.currentToken.Type == TokenElse || p.currentToken.Type == TokenElif {
			break
		}

		// Пропускаем пустые строки
		if p.currentToken.Type == TokenNewline {
			p.nextToken()
			continue
		}

		// Простая проверка: если встретили оператор на том же уровне, что и блок, выходим
		// Это упрощенная версия, в реальности нужно отслеживать отступы
		if p.currentToken.Type == TokenDef || p.currentToken.Type == TokenIf ||
			p.currentToken.Type == TokenFor || p.currentToken.Type == TokenWhile {
			// Проверяем, не является ли это вложенным блоком
			// Для простоты предполагаем, что это конец текущего блока
			break
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			statements = append(statements, stmt)
		}

		// Пропускаем newline после оператора
		if p.currentToken.Type == TokenNewline {
			p.nextToken()
		}

		// Для простоты: блок заканчивается после одного оператора
		// В реальной реализации нужно проверять отступы
		break
	}

	return statements, nil
}
