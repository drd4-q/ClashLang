# ClashLang Compiler Documentation

## Overview

The clashlang-compiler is a separate program that performs pre-processing, encryption, and security analysis before generating protected versions of ClashLang programs.

## Quick Start

Compile with memory protection:
```bash
./compiler myprogram.clash --protect
# Creates: myprogram.clash.secure.clang
```

Analyze for security issues:
```bash
./compiler myprogram.clash --analyze
```

Debug compilation details:
```bash
./compiler myprogram.clash --debug
```

## Critical Features

### 1. Memory Encryption
- All variables encrypted with XOR (key: 0x5A)
- Hex-encoded output for stealth
- Integrity hashes with SHA-256

### 2. Anti-Debugging
- OS-level thread locking
- Disables LD_PRELOAD/DYLD_INSERTIONS
- Reduces instrumentation visibility

### 3. Security Analysis
- Syntax validation with detailed error reporting
- Vulnerability detection:
  - Buffer overflow patterns
  - Use-after-free scenarios
  - Unsafe loop limits
  - Stack overflow risks

### 4. Integrity Verification
- Hash-based tamper detection
- Salted SHA-256 with runtime key
- Change detection on variables

## Integration with Main Interpreter

Current limitation: Protected binaries require custom interpreter patching to support ".secure.clang" extension.

## Usage Examples

### Basic Compilation:
```bash
# Generate protected binary
./compiler game.clash --protect

# Analyze for issues  
./compiler game.clash --analyze

# Debug compilation
./compiler game.clash --debug
```

### Integration with Main Interpreter:
```bash
# Modify clashlang binary if needed
# Add support for .secure.clang extension
# Run protected programs with enhanced security

# Or use original interpreter with analysis:
./clashlang program.clash --debug
```

## File Operations

### Protected Binary (.secure.clang):
- Encrypted variables with integrity hashes
- Embedded anti-debugging measures
- Detailed security analysis output
- Custom protection header

### Compatibility:
- Most ClashLang features supported
- Function registration maintained
- Memory operations enhanced
- Debugging enhanced with debug mode

## Build Instructions

```bash
# Build the compiler
Go build -o compiler compiler.go

# Generate protected binary
./compiler program.clash --protect

# Analysis mode
./compiler program.clash --analyze

# Debug mode
./compiler program.clash --debug
```

## Security Analysis Output

### Syntax Errors
Line-by-line error reporting for invalid syntax

### Security Warnings
Identifies potential vulnerabilities:
- Buffer overflow patterns
- Use-after-free scenarios
- Unsafe loop limits (0, 1 iterations)
- Stack overflow risks

### Protection Summary
Generated protection metrics:
- Number of protected variables
- Integrity hash status
- Anti-debugging activation
- Memory encryption confirmation
