package main

import (
	"fmt"
	"sync"
	"time"
)

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

type FunctionResult struct {
	result interface{}
	err    error
}

type Function func(int) (interface{}, error)

type Memory struct {
	f     Function
	cache map[int]FunctionResult
	lock  sync.RWMutex
}

func NewMemo(f Function) *Memory {
	return &Memory{
		f:     f,
		cache: map[int]FunctionResult{},
		lock:  sync.RWMutex{},
	}
}

func (m *Memory) Eval(key int) (interface{}, error) {
	m.lock.RLock()
	result, exists := m.cache[key]
	m.lock.RUnlock()
	if !exists {
		result.result, result.err = m.f(key)
		m.lock.Lock()
		m.cache[key] = result
		m.lock.Unlock()
	}
	return result.result, result.err
}

func FibonacciMemoWrapper(key int) (interface{}, error) {
	return Fibonacci(key), nil
}

func main() {
	fiboMemo := NewMemo(FibonacciMemoWrapper)
	fiboIndices := []int{40, 40, 40}
	var wg sync.WaitGroup
	for _, index := range fiboIndices {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			startFiboCalc := time.Now()
			result, err := fiboMemo.Eval(index)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Calculated Fibo(%d) = %d. It took %v\n", index, result, time.Since(startFiboCalc))
		}(index)
	}
	wg.Wait()
}
