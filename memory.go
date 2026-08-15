package main

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type Memory struct {
	variables map[string]interface{}
	functions map[string][]string
	key       byte
}

func NewMemory() *Memory {
	return &Memory{
		variables: make(map[string]interface{}),
		functions: make(map[string][]string),
		key:       0x5A,
	}
}

func (m *Memory) Set(name string, value interface{}) {
	m.variables[name] = value
}

func (m *Memory) Get(name string) (interface{}, bool) {
	val, ok := m.variables[name]
	return val, ok
}

func (m *Memory) StoreFunc(name string, lines []string) {
	m.functions[name] = lines
}

func (m *Memory) GetFunc(name string) ([]string, bool) {
	cmds, ok := m.functions[name]
	return cmds, ok
}

func (m *Memory) ClearVars() {
	m.variables = make(map[string]interface{})
}

func (m *Memory) xorEncode(s string) string {
	data := []byte(s)
	for i := range data {
		data[i] ^= m.key
	}
	return hex.EncodeToString(data)
}

func (m *Memory) xorDecode(hexStr string) string {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return hexStr
	}
	for i := range data {
		data[i] ^= m.key
	}
	return string(data)
}

func (m *Memory) encodeValue(val interface{}) string {
	s := fmt.Sprintf("%v", val)
	return m.xorEncode(s)
}

func (m *Memory) PrintAll() {
	keys := make([]string, 0, len(m.variables))
	for k := range m.variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for idx, key := range keys {
		encoded := m.encodeValue(m.variables[key])
		fmt.Printf("%d) %s:%s\n", idx+1, key, encoded)
	}
}

func (m *Memory) PrintDebug() {
	fmt.Println("=== DEBUG: Variables ===")
	keys := make([]string, 0, len(m.variables))
	for k := range m.variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for idx, key := range keys {
		fmt.Printf("%d) %s = %v\n", idx+1, key, m.variables[key])
	}

	fmt.Println("\n=== DEBUG: Functions ===")
	funcKeys := make([]string, 0, len(m.functions))
	for k := range m.functions {
		funcKeys = append(funcKeys, k)
	}
	sort.Strings(funcKeys)
	for idx, name := range funcKeys {
		lines := m.functions[name]
		fmt.Printf("%d) %s (%d lines)\n", idx+1, name, len(lines))
		for _, line := range lines {
			fmt.Printf("     %s\n", line)
		}
	}
}

func (m *Memory) DecodePrint() {
	fmt.Println("=== Decoded Memory ===")
	keys := make([]string, 0, len(m.variables))
	for k := range m.variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for idx, key := range keys {
		fmt.Printf("%d) %s = %v\n", idx+1, key, m.variables[key])
	}
}

func isHexString(s string) bool {
	if len(s) == 0 || len(s)%2 != 0 {
		return false
	}
	for _, c := range s {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}
