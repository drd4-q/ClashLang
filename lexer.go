package main

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType определяет тип токена
type TokenType int

const (
	// Специальные токены
	TokenEOF TokenType = iota
	TokenNewline
	TokenIndent
	TokenDedent

	// Литералы
	TokenNumber
	TokenString
	TokenIdentifier

	// Ключевые слова
	TokenPrint
	TokenIf
	TokenElse
	TokenElif
	TokenFor
	TokenWhile
	TokenDef
	TokenReturn
	TokenBreak
	TokenContinue
	TokenIn
	TokenRange

	// Управление памятью
	TokenMemAlloc
	TokenMemFree
	TokenMemCheck

	// Управление процессами
	TokenProcCreate
	TokenProcRead
	TokenProcKill
	TokenProcList

	// Операторы
	TokenAssign    // =
	TokenPlus      // +
	TokenMinus     // -
	TokenMultiply  // *
	TokenDivide    // /
	TokenModulo    // %
	TokenPower     // **
	TokenEqual     // ==
	TokenNotEqual  // !=
	TokenLess      // <
	TokenGreater   // >
	TokenLessEq    // <=
	TokenGreaterEq // >=
	TokenAnd       // and
	TokenOr        // or
	TokenNot       // not

	// Разделители
	TokenLeftParen  // (
	TokenRightParen // )
	TokenLeftBrace  // [
	TokenRightBrace // ]
	TokenLeftCurly  // {
	TokenRightCurly // }
	TokenComma      // ,
	TokenColon      // :
	TokenDot        // .
)

// Token представляет токен
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// Lexer выполняет лексический анализ
type Lexer struct {
	input        string
	position     int
	line         int
	column       int
	indentStack  []int
	pendingTokens []Token
}

// NewLexer создает новый лексер
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:        input,
		position:     0,
		line:         1,
		column:       1,
		indentStack:  []int{0},
		pendingTokens: []Token{},
	}
}

var keywords = map[string]TokenType{
	"print":       TokenPrint,
	"if":          TokenIf,
	"else":        TokenElse,
	"elif":        TokenElif,
	"for":         TokenFor,
	"while":       TokenWhile,
	"def":         TokenDef,
	"return":      TokenReturn,
	"break":       TokenBreak,
	"continue":    TokenContinue,
	"in":          TokenIn,
	"range":       TokenRange,
	"and":         TokenAnd,
	"or":          TokenOr,
	"not":         TokenNot,
	"mem_alloc":   TokenMemAlloc,
	"mem_free":    TokenMemFree,
	"mem_check":   TokenMemCheck,
	"proc_create": TokenProcCreate,
	"proc_read":   TokenProcRead,
	"proc_kill":   TokenProcKill,
	"proc_list":   TokenProcList,
}

// NextToken возвращает следующий токен
func (l *Lexer) NextToken() Token {
	// Если есть отложенные токены, вернуть их
	if len(l.pendingTokens) > 0 {
		token := l.pendingTokens[0]
		l.pendingTokens = l.pendingTokens[1:]
		return token
	}

	l.skipWhitespace()

	if l.position >= len(l.input) {
		// Генерируем DEDENT токены для закрытия всех уровней отступов
		if len(l.indentStack) > 1 {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			return Token{Type: TokenDedent, Line: l.line, Column: l.column}
		}
		return Token{Type: TokenEOF, Line: l.line, Column: l.column}
	}

	ch := l.input[l.position]

	// Комментарии
	if ch == '#' {
		l.skipComment()
		return l.NextToken()
	}

	// Строки
	if ch == '"' || ch == '\'' {
		return l.readString(ch)
	}

	// Числа
	if unicode.IsDigit(rune(ch)) {
		return l.readNumber()
	}

	// Идентификаторы и ключевые слова
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		return l.readIdentifier()
	}

	// Операторы и разделители
	token := Token{Line: l.line, Column: l.column}

	switch ch {
	case '=':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			token.Type = TokenEqual
			token.Value = "=="
		} else {
			l.advance()
			token.Type = TokenAssign
			token.Value = "="
		}
	case '+':
		l.advance()
		token.Type = TokenPlus
		token.Value = "+"
	case '-':
		l.advance()
		token.Type = TokenMinus
		token.Value = "-"
	case '*':
		if l.peek() == '*' {
			l.advance()
			l.advance()
			token.Type = TokenPower
			token.Value = "**"
		} else {
			l.advance()
			token.Type = TokenMultiply
			token.Value = "*"
		}
	case '/':
		l.advance()
		token.Type = TokenDivide
		token.Value = "/"
	case '%':
		l.advance()
		token.Type = TokenModulo
		token.Value = "%"
	case '!':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			token.Type = TokenNotEqual
			token.Value = "!="
		}
	case '<':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			token.Type = TokenLessEq
			token.Value = "<="
		} else {
			l.advance()
			token.Type = TokenLess
			token.Value = "<"
		}
	case '>':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			token.Type = TokenGreaterEq
			token.Value = ">="
		} else {
			l.advance()
			token.Type = TokenGreater
			token.Value = ">"
		}
	case '(':
		l.advance()
		token.Type = TokenLeftParen
		token.Value = "("
	case ')':
		l.advance()
		token.Type = TokenRightParen
		token.Value = ")"
	case '[':
		l.advance()
		token.Type = TokenLeftBrace
		token.Value = "["
	case ']':
		l.advance()
		token.Type = TokenRightBrace
		token.Value = "]"
	case '{':
		l.advance()
		token.Type = TokenLeftCurly
		token.Value = "{"
	case '}':
		l.advance()
		token.Type = TokenRightCurly
		token.Value = "}"
	case ',':
		l.advance()
		token.Type = TokenComma
		token.Value = ","
	case ':':
		l.advance()
		token.Type = TokenColon
		token.Value = ":"
	case '.':
		l.advance()
		token.Type = TokenDot
		token.Value = "."
	case '\n':
		l.advance()
		l.line++
		l.column = 1
		token.Type = TokenNewline
		token.Value = "\\n"
	default:
		l.advance()
		token.Type = TokenEOF
		token.Value = string(ch)
	}

	return token
}

func (l *Lexer) advance() {
	if l.position < len(l.input) {
		l.position++
		l.column++
	}
}

func (l *Lexer) peek() byte {
	if l.position+1 < len(l.input) {
		return l.input[l.position+1]
	}
	return 0
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) {
		ch := l.input[l.position]
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *Lexer) skipComment() {
	for l.position < len(l.input) && l.input[l.position] != '\n' {
		l.advance()
	}
}

func (l *Lexer) readString(quote byte) Token {
	token := Token{Type: TokenString, Line: l.line, Column: l.column}
	l.advance() // Пропускаем открывающую кавычку

	var value strings.Builder
	for l.position < len(l.input) {
		ch := l.input[l.position]
		if ch == quote {
			l.advance() // Пропускаем закрывающую кавычку
			break
		}
		if ch == '\\' && l.position+1 < len(l.input) {
			l.advance()
			next := l.input[l.position]
			switch next {
			case 'n':
				value.WriteByte('\n')
			case 't':
				value.WriteByte('\t')
			case '\\':
				value.WriteByte('\\')
			case quote:
				value.WriteByte(quote)
			default:
				value.WriteByte(next)
			}
			l.advance()
		} else {
			value.WriteByte(ch)
			l.advance()
		}
	}

	token.Value = value.String()
	return token
}

func (l *Lexer) readNumber() Token {
	token := Token{Type: TokenNumber, Line: l.line, Column: l.column}
	var value strings.Builder

	for l.position < len(l.input) {
		ch := l.input[l.position]
		if unicode.IsDigit(rune(ch)) || ch == '.' {
			value.WriteByte(ch)
			l.advance()
		} else {
			break
		}
	}

	token.Value = value.String()
	return token
}

func (l *Lexer) readIdentifier() Token {
	token := Token{Type: TokenIdentifier, Line: l.line, Column: l.column}
	var value strings.Builder

	for l.position < len(l.input) {
		ch := l.input[l.position]
		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
			value.WriteByte(ch)
			l.advance()
		} else {
			break
		}
	}

	token.Value = value.String()

	// Проверяем, является ли это ключевым словом
	if tokenType, isKeyword := keywords[token.Value]; isKeyword {
		token.Type = tokenType
	}

	return token
}

// TokenTypeString возвращает строковое представление типа токена
func TokenTypeString(t TokenType) string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNewline:
		return "NEWLINE"
	case TokenIndent:
		return "INDENT"
	case TokenDedent:
		return "DEDENT"
	case TokenNumber:
		return "NUMBER"
	case TokenString:
		return "STRING"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenPrint:
		return "PRINT"
	case TokenIf:
		return "IF"
	case TokenElse:
		return "ELSE"
	case TokenElif:
		return "ELIF"
	case TokenFor:
		return "FOR"
	case TokenWhile:
		return "WHILE"
	case TokenDef:
		return "DEF"
	case TokenReturn:
		return "RETURN"
	case TokenBreak:
		return "BREAK"
	case TokenContinue:
		return "CONTINUE"
	case TokenIn:
		return "IN"
	case TokenRange:
		return "RANGE"
	case TokenMemAlloc:
		return "MEM_ALLOC"
	case TokenMemFree:
		return "MEM_FREE"
	case TokenMemCheck:
		return "MEM_CHECK"
	case TokenProcCreate:
		return "PROC_CREATE"
	case TokenProcRead:
		return "PROC_READ"
	case TokenProcKill:
		return "PROC_KILL"
	case TokenProcList:
		return "PROC_LIST"
	case TokenAssign:
		return "ASSIGN"
	case TokenPlus:
		return "PLUS"
	case TokenMinus:
		return "MINUS"
	case TokenMultiply:
		return "MULTIPLY"
	case TokenDivide:
		return "DIVIDE"
	case TokenModulo:
		return "MODULO"
	case TokenPower:
		return "POWER"
	case TokenEqual:
		return "EQUAL"
	case TokenNotEqual:
		return "NOT_EQUAL"
	case TokenLess:
		return "LESS"
	case TokenGreater:
		return "GREATER"
	case TokenLessEq:
		return "LESS_EQ"
	case TokenGreaterEq:
		return "GREATER_EQ"
	case TokenAnd:
		return "AND"
	case TokenOr:
		return "OR"
	case TokenNot:
		return "NOT"
	case TokenLeftParen:
		return "LEFT_PAREN"
	case TokenRightParen:
		return "RIGHT_PAREN"
	case TokenLeftBrace:
		return "LEFT_BRACE"
	case TokenRightBrace:
		return "RIGHT_BRACE"
	case TokenLeftCurly:
		return "LEFT_CURLY"
	case TokenRightCurly:
		return "RIGHT_CURLY"
	case TokenComma:
		return "COMMA"
	case TokenColon:
		return "COLON"
	case TokenDot:
		return "DOT"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}
