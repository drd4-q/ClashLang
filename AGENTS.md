# ClashLang

Go-based interpreter for ClashLang with function storage, variable management, and file operations.

## Build & Run

```bash
go build -o clashlang . && ./clashlang program.clash
./clashlang program.clash --debug  # with debug output
```

## Critical Security Tool

**Enhanced Security Compiler** (`compiler.go`):
- Memory encryption and protection
- Security analysis and vulnerability detection
- Anti-debugging measures
- Integrity verification

**Usage**:
```bash
./compiler program.clash --protect  # Creates protected .secure.clang
./compiler program.clash --analyze  # Analyzes for security issues
./compiler program.clash --debug    # Detailed compilation info
```

## Architecture

- `main.go` — entry point, validates .clash files
- `interpreter.go` — core execution logic
- `memory.go` — XOR encryption engine
- `file_ops.go` — file operations
- `compiler.go` — extended protection (separate tool)
- `clashlang` binary — main interpreter executable

## Critical Operations (Likely to be missed)

**Function Registration**:
1. Full definition: `function(Print){...}`
2. OR bare: `function(MyFunc)` + `memory.start{...}`
3. Execution: `memory.load()` or `memory.start()`

**Example Program**:
```clash
function(c){  // With body
c = 12
print(c)
}

memory.start{    // Separate body
	memory.start(c)
	c = 25
}

memory.load(c)
memory.out  // Encrypted: c:7A...
```

## Syntax Rules (case-insensitive)

**Variable Assignment**:
```clash
c = 42        // int
pi = 3.14     // float  
name = 'John' // string
flag = true   // bool
```

**Commands**:
- Case-insensitive: `print()`, `solve()`, `jump()`
- Case-sensitive: `x` vs `X` (variables)
- Two-line ifs: `if x = 12` / `then jump(Func)`

**Memory**:
- `memory.out` (encrypted)
- `debug()` (plaintext)
- `decode()` (decrypt)

**File Safety**:
- `file.open(name)` / `file.close`
- `file.text.read` / `.copy = var`
- `file.create(name)`

**String Operations**:
- `len(var)`, `upper()`, `lower()`, `reverse()`
- `replace(var, old, new)`
- `substr(var, start, len)`

**Math**:
- `solve(a + b)` → `solve.out = var`
- Operators: `+`, `-`, `*`, `/`, `^`

**Control Flow**:
- `if expr then jump(name) else: skip`
- `jump(name)` (calls functions)
- `memory.start(name)` (alternate call)
- `repeat(n){...}` (loops)

## Protection Features

**XOR Encryption**:
- Key: 0x5A
- All variables encrypted in `memory.out`
- Hex-encoded for stealth

**Security Analysis**:
- Syntax validation
- Buffer overflow detection
- Use-after-free detection
- Stack overflow warnings

**Anti-Debugging**:
- OS-level thread locking
- Disables dynamic linking
- Memory tampering detection

**Memory Safety**:
- Integrity hashes (SHA-256)
- Tamper verification
- Runtime protection

## Basic Workflow

1. Write `.clash` file with functions
2. Test with analysis:
   ```bash
   ./compiler program.clash --analyze
   ```
3. Generate secure version:
   ```bash
   ./compiler program.clash --protect
   ```
4. Execute:
   ```bash
   ./clashlang program.clash --debug
   ```

## File Requirements

- Files must end with `.clash` extension
- Comments: `// comment`
- No tests/linting/CI configured
