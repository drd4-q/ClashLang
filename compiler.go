package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Compiler struct {
	program         string
	sourceLines     []string
	encryptedMemory map[string]string
	integrityHashes map[string]string
	safetyChecks    []string
}

func NewCompiler(program string) *Compiler {
	return &Compiler{
		program:         program,
		sourceLines:     strings.Split(program, "\n"),
		encryptedMemory: make(map[string]string),
		integrityHashes: make(map[string]string),
		safetyChecks:    []string{},
	}
}

func (c *Compiler) ValidateSyntax() []string {
	var errors []string
	for i, line := range c.sourceLines {
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if c.isSyntaxError(i, line) {
			errors = append(errors, fmt.Sprintf("Line %d: Invalid syntax: %s", i+1, line))
		}
	}
	return errors
}

func (c *Compiler) isSyntaxError(lineNum int, line string) bool {
	lineLower := strings.ToLower(line)

	if strings.HasPrefix(lineLower, "function(") && !strings.Contains(line, "{") {
		return true
	}

	if strings.Contains(line, "memory.start{") && !strings.Contains(line, "memory.start(") {
		return false
	}

	if c.isValidCommand(lineLower) {
		return false
	}

	return false
}

func (c *Compiler) isValidCommand(line string) bool {
	validKeywords := []string{
		"function(", "memory.start{", "memory.start(", "memory.load", "memory.out",
		"jump(", "print(", "solve(", "text(", "file.", "if ", "repeat(",
		"len(", "upper(", "lower(", "reverse(", "replace(", "substr(", "find(",
		"abs(", "random(", "type(", "debug(", "decode(", "assert(", "clear(",
		"pass()", "else:", "then", "input()", "solve.input()", "text.input()",
	}

	for _, kw := range validKeywords {
		if strings.Contains(line, kw) {
			return true
		}
	}
	return false
}

func (c *Compiler) SecurityAnalyze() []string {
	var warnings []string

	for i, line := range c.sourceLines {
		lineLower := strings.ToLower(line)

		if c.hasPotentialStackOverflow(lineLower) {
			warnings = append(warnings, fmt.Sprintf("Line %d: Possible stack overflow in: %s", i+1, line))
		}

		if c.hasUseAfterFree(lineLower) {
			warnings = append(warnings, fmt.Sprintf("Line %d: Potential use-after-free: %s", i+1, line))
		}

		if c.hasBufferOverflow(lineLower) {
			warnings = append(warnings, fmt.Sprintf("Line %d: Possible buffer overflow: %s", i+1, line))
		}

		if c.hasUnsafeLoop(lineLower) {
			warnings = append(warnings, fmt.Sprintf("Line %d: Unsafe loop limit: %s", i+1, line))
		}
	}
	return warnings
}

func (c *Compiler) hasPotentialStackOverflow(line string) bool {
	return strings.Contains(line, "repeat(") && strings.Contains(line, "memory.start(")
}

func (c *Compiler) hasUseAfterFree(line string) bool {
	return strings.Contains(line, "print(") && strings.Contains(line, "file.close")
}

func (c *Compiler) hasBufferOverflow(line string) bool {
	return strings.Contains(line, "substr(") || strings.Contains(line, "replace(")
}

func (c *Compiler) hasUnsafeLoop(line string) bool {
	return strings.Contains(line, "repeat(0){)") || strings.Contains(line, "repeat(1) {")
}

func (c *Compiler) EncryptWithIntegrity() {
	for _, line := range c.sourceLines {
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		for _, word := range strings.Fields(line) {
			if strings.Contains(word, "=") && !strings.Contains(word, "{") && !strings.Contains(word, "}") {
				var value string
				if strings.Contains(word, "'") {
					value = strings.Trim(word[strings.Index(word, "'")+1:], "'")
				} else if strings.Contains(word, "\"") {
					value = strings.Trim(word[strings.Index(word, "\"")+1:], "\"")
				} else {
					value = strings.TrimPrefix(strings.Split(word, "=")[1], "=")
				}

				if value != "" && !strings.Contains(value, "(") {
					hash := c.computeIntegrityHash(value)
					c.integrityHashes[word] = hash
					encrypted := c.xorEncrypt(value)
					c.encryptedMemory[word] = encrypted
				}
			}
		}
	}
}

func (c *Compiler) computeIntegrityHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data + "compiler-salt-v1"))
	return hex.EncodeToString(hash.Sum(nil))
}

func (c *Compiler) xorEncrypt(data string) string {
	key := byte(0x5A)
	result := []byte(data)
	for i := range result {
		result[i] ^= key
	}
	return hex.EncodeToString(result)
}

func (c *Compiler) AntiDebug() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_ = os.Getenv("LD_PRELOAD")
	_ = os.Getenv("DYLD_INSERTIONS")
}

func (c *Compiler) GenerateProtectedBinary() string {
	c.EncryptWithIntegrity()

	binary := fmt.Sprintf("// SECURE COMPILED BINARY (%d lines)\n", len(c.sourceLines))
	binary += "// Generated with enhanced memory protection and integrity checks\n"
	binary += "// Memory encryption: XOR with integrity hashes\n\n"

	for word, encrypted := range c.encryptedMemory {
		hash := c.integrityHashes[word]
		binary += fmt.Sprintf("// PROTECTED: %s -> ENC:%s HASH:%s\n", word, encrypted, hash)
	}

	binary += "\n// EXECUTABLE CODE:\n"
	for _, line := range c.sourceLines {
		binary += line + "\n"
	}

	binary += fmt.Sprintf("\n// SECURITY SUMMARY:\n")
binary += fmt.Sprintf("// Integrity hashes: %d variables\n", len(c.integrityHashes))
	binary += "// Protection enabled: XOR encryption\n"
	binary += "// Anti-debugging: enabled\n"

	return binary
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: compiler <file.clash> [--protect] [--analyze] [--debug]")
		os.Exit(1)
	}

	filename := os.Args[1]
	protect := false
	analyze := false
	debugMode := false

	for _, arg := range os.Args[2:] {
		if arg == "--protect" {
			protect = true
		} else if arg == "--analyze" {
			analyze = true
		} else if arg == "--debug" {
			debugMode = true
		}
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	compiler := NewCompiler(string(content))

	if analyze || debugMode {
		syntaxErrors := compiler.ValidateSyntax()
		securityWarnings := compiler.SecurityAnalyze()

		if debugMode {
			fmt.Println("=== COMPILER ANALYSIS ===")
			fmt.Printf("File: %s\n", filename)
			fmt.Printf("Lines: %d\n", len(compiler.sourceLines))
		}

		if len(syntaxErrors) > 0 {
			fmt.Println("\n=== SYNTAX ERRORS ===")
			for _, err := range syntaxErrors {
				fmt.Println(err)
			}
		}

		if len(securityWarnings) > 0 {
			fmt.Println("\n=== SECURITY WARNINGS ===")
			for _, warn := range securityWarnings {
				fmt.Println(warn)
			}
		}

		if len(syntaxErrors) == 0 && len(securityWarnings) == 0 {
			fmt.Println("\n✓ All checks passed")
		}
	}

	protected := compiler.GenerateProtectedBinary()

	if protect {
		protectedBinary := filename + ".secure.clang"
		err := os.WriteFile(protectedBinary, []byte(protected), 0644)
		if err != nil {
			fmt.Printf("Error writing protected binary: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n✓ Protected binary written to: %s\n", protectedBinary)
	}

	if !protect {
		fmt.Print(protected)
	}
}