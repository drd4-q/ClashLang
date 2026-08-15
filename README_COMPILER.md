ClashLang Compiler v1.0

Advanced compiler with memory protection and security analysis for ClashLang programs.

## Features

### Memory Protection
- **XOR Encryption**: All variables encrypted with key 0x5A
- **Integrity Hashes**: SHA-256 hashes for each variable
- **Anti-Debugging**: Disables debugging tools at OS level
- **Memory Tampering Detection**: Detects unauthorized memory modifications

### Security Analysis
- **Syntax Validation**: Catches invalid function registrations and syntax errors
- **Vulnerability Detection**: Identifies potential security issues:
  - Buffer overflow attempts
  - Use-after-free scenarios
  - Unsafe loop limits
  - Stack overflow risks
- **Anti-Static Analysis**: Uses memory-randomized keys

### Integration Modes

#### Standard Mode
```bash
./compiler program.clash
```
Generates protected .secure.clang binary with encrypted variables

#### Analysis Mode
```bash
./compiler program.clash --analyze
./compiler program.clash --debug
```
Validates syntax and reports security warnings

#### Combined Mode
```bash
./compiler program.clash --protect --analyze --debug
```
Full protection with security verification

## Architecture

The compiler transforms CliskLang source into secure binaries by:

1. **Pre-processing**: Validates syntax and detects vulnerabilities
2. **Encryption**: Encrypts all variables with XOR
3. **Hashing**: Computes integrity hashes for verification
4. **Protection**: Applies anti-debugging measures
5. **Generation**: Creates .secure.clang with embedded protections

## Security Notes

- **Warning**: This compiler is "bad in security" intentionally - it creates over-protective binaries
- **Vulnerability Simulation**: Designed to demonstrate potential security issues
- **Debugging Features**: Reduces obfuscating to aid educational purposes

## Usage Examples

### Compile with protection:
```bash
./compiler myprogram.clash --protect
# Creates: myprogram.clash.secure.clang
```

### Analyze for issues:
```bash
./compiler myprogram.clash --analyze
# Reports syntax errors and security warnings
```

### Debug compilation:
```bash
./compiler myprogram.clash --debug
# Shows detailed analysis of program structure
```

## Integration with Main Interpreter

To use protected binaries:
1. Compile with `compiler.py --protect`
2. Edit original clashlang binary to use protected extension
3. Execute normal interpreter