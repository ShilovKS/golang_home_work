package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return nil // Или ошибка, если n должно быть > 0
	}

	taskCh := make(chan Task)
	stop := make(chan struct{})
	var wg sync.WaitGroup
	var errCount int32

	// Запуск воркеров
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if err := task(); err != nil {
					if m > 0 {
						if atomic.AddInt32(&errCount, 1) == int32(m) {
							close(stop)
						}
					}
				}
			}
		}()
	}

	// Отправка задач с учётом остановки
	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			select {
			case taskCh <- task:
			case <-stop:
				return
			}
		}
	}()

	// Ожидание завершения
	select {
	case <-stop:
		return ErrErrorsLimitExceeded
	case <-wait(&wg):
		if m > 0 && atomic.LoadInt32(&errCount) >= int32(m) {
			return ErrErrorsLimitExceeded
		}
		return nil
	}
}

func wait(wg *sync.WaitGroup) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	return done
}
