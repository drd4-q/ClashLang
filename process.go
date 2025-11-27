package main

import (
	"fmt"
	"sync"
	"time"
)

// ProcessState определяет состояние процесса
type ProcessState int

const (
	ProcessRunning ProcessState = iota
	ProcessStopped
	ProcessFinished
)

func (ps ProcessState) String() string {
	switch ps {
	case ProcessRunning:
		return "Running"
	case ProcessStopped:
		return "Stopped"
	case ProcessFinished:
		return "Finished"
	default:
		return "Unknown"
	}
}

// Process представляет процесс
type Process struct {
	ID         int
	Name       string
	Code       string
	State      ProcessState
	CreateTime time.Time
	Variables  map[string]interface{}
}

// ProcessManager управляет процессами
type ProcessManager struct {
	processes  map[int]*Process
	nextID     int
	mutex      sync.RWMutex
}

// NewProcessManager создает новый менеджер процессов
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		processes: make(map[int]*Process),
		nextID:    1,
	}
}

// Create создает новый процесс
func (pm *ProcessManager) Create(name string, code string) int {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	id := pm.nextID
	pm.nextID++

	process := &Process{
		ID:         id,
		Name:       name,
		Code:       code,
		State:      ProcessRunning,
		CreateTime: time.Now(),
		Variables:  make(map[string]interface{}),
	}

	pm.processes[id] = process

	fmt.Printf("[PROCESS] Создан процесс #%d '%s'\n", id, name)
	return id
}

// Read читает код процесса
func (pm *ProcessManager) Read(id int) (string, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	process, exists := pm.processes[id]
	if !exists {
		return "", fmt.Errorf("процесс #%d не найден", id)
	}

	return process.Code, nil
}

// Kill завершает процесс
func (pm *ProcessManager) Kill(id int) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	process, exists := pm.processes[id]
	if !exists {
		return fmt.Errorf("процесс #%d не найден", id)
	}

	process.State = ProcessFinished
	fmt.Printf("[PROCESS] Процесс #%d '%s' завершен\n", id, process.Name)
	return nil
}

// Get получает процесс по ID
func (pm *ProcessManager) Get(id int) (*Process, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	process, exists := pm.processes[id]
	if !exists {
		return nil, fmt.Errorf("процесс #%d не найден", id)
	}

	return process, nil
}

// List выводит список всех процессов
func (pm *ProcessManager) List() {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	if len(pm.processes) == 0 {
		fmt.Println("\n=== СПИСОК ПРОЦЕССОВ ===")
		fmt.Println("Нет активных процессов")
		fmt.Println("========================\n")
		return
	}

	fmt.Println("\n=== СПИСОК ПРОЦЕССОВ ===")
	fmt.Printf("Всего процессов: %d\n\n", len(pm.processes))

	for id, proc := range pm.processes {
		duration := time.Since(proc.CreateTime)
		fmt.Printf("Процесс #%d:\n", id)
		fmt.Printf("  Имя: %s\n", proc.Name)
		fmt.Printf("  Состояние: %s\n", proc.State)
		fmt.Printf("  Время работы: %v\n", duration.Round(time.Millisecond))
		fmt.Printf("  Переменных: %d\n", len(proc.Variables))
		fmt.Println()
	}
	fmt.Println("========================\n")
}

// SetVariable устанавливает переменную процесса
func (pm *ProcessManager) SetVariable(id int, name string, value interface{}) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	process, exists := pm.processes[id]
	if !exists {
		return fmt.Errorf("процесс #%d не найден", id)
	}

	process.Variables[name] = value
	return nil
}

// GetVariable получает переменную процесса
func (pm *ProcessManager) GetVariable(id int, name string) (interface{}, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	process, exists := pm.processes[id]
	if !exists {
		return nil, fmt.Errorf("процесс #%d не найден", id)
	}

	value, exists := process.Variables[name]
	if !exists {
		return nil, fmt.Errorf("переменная '%s' не найдена в процессе #%d", name, id)
	}

	return value, nil
}

// GetCount возвращает количество процессов
func (pm *ProcessManager) GetCount() int {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	return len(pm.processes)
}

// Clear очищает все процессы (для тестирования)
func (pm *ProcessManager) Clear() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.processes = make(map[int]*Process)
	pm.nextID = 1
}
