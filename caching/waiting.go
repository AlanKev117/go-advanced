package main

import (
	"fmt"
	"sync"
	"time"
)

// This function represents a highly expensive-to-run function.
// Replace the body with the required behavior
func ExpensiveFunction(n interface{}) interface{} {
	fmt.Printf("Executing expensive calculation for %v\n", n)
	time.Sleep(5 * time.Second)
	return n
}

// This struct is meant to wrap a function to make it
// concurrent-safe and to cache its results in memory
type Service struct {
	InProgress     map[interface{}]bool
	PendingReaders map[interface{}][]chan interface{}
	Lock           sync.RWMutex
	f              func(interface{}) interface{}
}

// Excecutes the wrapped function in Service s and caches the result for further calls.
// If called concurrently with same arguments, it puts newer processes to wait until
// first call finishes, following DRY principle.
func (s *Service) Exec(job interface{}) {

	s.Lock.RLock()
	// In case this routine has to wait for the calculation routine:
	if s.InProgress[job] {

		// Make a new channel (reader)
		s.Lock.RUnlock()
		response := make(chan interface{})
		defer close(response)

		// Append the reader to its respective queue (regarding to the job)
		s.Lock.Lock()
		s.PendingReaders[job] = append(s.PendingReaders[job], response)
		s.Lock.Unlock()

		// The reader now only has to wait for the result of the routine
		// that is calculating the result for the required job
		fmt.Printf("Waiting for Response job: %d\n", job)
		resp := <-response
		fmt.Printf("Response Done, received %d\n", resp)
		return
	}

	// In case this routine is the one to perform the calculation:

	s.Lock.RUnlock()
	s.Lock.Lock()

	// Flag to tell other routines that the required job is being calculated
	s.InProgress[job] = true
	s.Lock.Unlock()

	fmt.Printf("Performing expensive function for job %d\n", job)
	result := s.f(job)

	// Once finished the function call, recall the channels
	// to send the result to
	s.Lock.RLock()
	pendingWorkers, inProgress := s.PendingReaders[job]
	s.Lock.RUnlock()

	// Send the message to all routines via channels
	if inProgress {
		for _, pendingWorker := range pendingWorkers {
			pendingWorker <- result
		}
		fmt.Printf("Result sent - all pending workers ready job:%d\n", job)
	}

	// Free in-progess flag and delete all readers
	s.Lock.Lock()
	s.InProgress[job] = false
	s.PendingReaders[job] = make([]chan interface{}, 0)
	s.Lock.Unlock()
}

// Constructor function for a Service object
func NewService(f func(interface{}) interface{}) *Service {
	return &Service{
		InProgress:     make(map[interface{}]bool),
		PendingReaders: make(map[interface{}][]chan interface{}),
		f:              f,
	}
}

func main() {

	service := NewService(ExpensiveFunction)

	jobs := []interface{}{40, 40, 40}

	var wg sync.WaitGroup
	wg.Add(len(jobs))

	// Invoke as many routines as jobs are there
	for _, job := range jobs {
		go func(job interface{}) {
			defer wg.Done()
			// Same service will execute concurrently
			service.Exec(job)
		}(job)
	}

	wg.Wait()
}
