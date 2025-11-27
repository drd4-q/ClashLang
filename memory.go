package main

import (
	"fmt"
	"sync"
	"time"
)

// MemoryBlock представляет блок выделенной памяти
type MemoryBlock struct {
	Name      string
	Size      int
	Data      []byte
	AllocTime time.Time
}

// MemoryManager управляет памятью
type MemoryManager struct {
	blocks      map[string]*MemoryBlock
	totalAlloc  int
	totalFreed  int
	mutex       sync.RWMutex
}

// NewMemoryManager создает новый менеджер памяти
func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		blocks:     make(map[string]*MemoryBlock),
		totalAlloc: 0,
		totalFreed: 0,
	}
}

// Alloc выделяет память
func (mm *MemoryManager) Alloc(name string, size int) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	if _, exists := mm.blocks[name]; exists {
		return fmt.Errorf("память с именем '%s' уже выделена", name)
	}

	if size <= 0 {
		return fmt.Errorf("размер памяти должен быть положительным числом")
	}

	block := &MemoryBlock{
		Name:      name,
		Size:      size,
		Data:      make([]byte, size),
		AllocTime: time.Now(),
	}

	mm.blocks[name] = block
	mm.totalAlloc += size

	fmt.Printf("[MEMORY] Выделено %d байт для '%s'\n", size, name)
	return nil
}

// Free освобождает память
func (mm *MemoryManager) Free(name string) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	block, exists := mm.blocks[name]
	if !exists {
		return fmt.Errorf("память с именем '%s' не найдена", name)
	}

	mm.totalFreed += block.Size
	delete(mm.blocks, name)

	fmt.Printf("[MEMORY] Освобождено %d байт для '%s'\n", block.Size, name)
	return nil
}

// Get получает блок памяти
func (mm *MemoryManager) Get(name string) (*MemoryBlock, error) {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	block, exists := mm.blocks[name]
	if !exists {
		return nil, fmt.Errorf("память с именем '%s' не найдена", name)
	}

	return block, nil
}

// Write записывает данные в память
func (mm *MemoryManager) Write(name string, offset int, data []byte) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	block, exists := mm.blocks[name]
	if !exists {
		return fmt.Errorf("память с именем '%s' не найдена", name)
	}

	if offset < 0 || offset+len(data) > block.Size {
		return fmt.Errorf("выход за границы памяти")
	}

	copy(block.Data[offset:], data)
	return nil
}

// Read читает данные из памяти
func (mm *MemoryManager) Read(name string, offset int, length int) ([]byte, error) {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	block, exists := mm.blocks[name]
	if !exists {
		return nil, fmt.Errorf("память с именем '%s' не найдена", name)
	}

	if offset < 0 || offset+length > block.Size {
		return nil, fmt.Errorf("выход за границы памяти")
	}

	data := make([]byte, length)
	copy(data, block.Data[offset:offset+length])
	return data, nil
}

// Check проверяет состояние памяти и выводит статистику
func (mm *MemoryManager) Check() {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	fmt.Println("\n=== СТАТИСТИКА ПАМЯТИ ===")
	fmt.Printf("Всего выделено: %d байт\n", mm.totalAlloc)
	fmt.Printf("Всего освобождено: %d байт\n", mm.totalFreed)
	fmt.Printf("Текущее использование: %d байт\n", mm.totalAlloc-mm.totalFreed)
	fmt.Printf("Активных блоков: %d\n\n", len(mm.blocks))

	if len(mm.blocks) > 0 {
		fmt.Println("Активные блоки памяти:")
		for name, block := range mm.blocks {
			duration := time.Since(block.AllocTime)
			fmt.Printf("  - %s: %d байт (выделено %v назад)\n", name, block.Size, duration.Round(time.Millisecond))
		}
	}

	// Проверка утечек памяти
	if len(mm.blocks) > 0 {
		fmt.Println("\n⚠️  ВНИМАНИЕ: Обнаружены неосвобожденные блоки памяти!")
		fmt.Println("Рекомендуется освободить память с помощью mem_free()")
	} else {
		fmt.Println("\n✓ Утечек памяти не обнаружено")
	}
	fmt.Println("========================\n")
}

// GetStats возвращает статистику памяти
func (mm *MemoryManager) GetStats() (totalAlloc, totalFreed, current, blocks int) {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	return mm.totalAlloc, mm.totalFreed, mm.totalAlloc - mm.totalFreed, len(mm.blocks)
}

// Clear очищает всю память (для тестирования)
func (mm *MemoryManager) Clear() {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	mm.blocks = make(map[string]*MemoryBlock)
	mm.totalAlloc = 0
	mm.totalFreed = 0
}
