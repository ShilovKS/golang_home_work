package hw05parallelexecution

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

// TestRunErrorsLimit проверяет, что если в первых M задачах произошли ошибки,
// то общее число выполненных задач не превышает N+M.
func TestRunErrorsLimit(t *testing.T) {
	defer goleak.VerifyNone(t)

	tasksCount := 50
	tasks := make([]Task, 0, tasksCount)
	var runTasksCount int32

	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		tasks = append(tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
	}

	workersCount := 10
	maxErrorsCount := 23
	err := Run(tasks, workersCount, maxErrorsCount)

	require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual error - %v", err)
	// Используем атомарное чтение, чтобы избежать гонки:
	require.LessOrEqual(t, atomic.LoadInt32(&runTasksCount), int32(workersCount+maxErrorsCount), "extra tasks were started")
}

// TestRunNoErrors проверяет, что если все задачи выполняются без ошибок,
// то все они выполняются параллельно (общая длительность меньше суммы задержек).
func TestRunNoErrors(t *testing.T) {
	defer goleak.VerifyNone(t)

	tasksCount := 50
	tasks := make([]Task, 0, tasksCount)
	var runTasksCount int32
	var sumSleepTime time.Duration

	for i := 0; i < tasksCount; i++ {
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
		sumSleepTime += taskSleep

		tasks = append(tasks, func() error {
			time.Sleep(taskSleep)
			atomic.AddInt32(&runTasksCount, 1)
			return nil
		})
	}

	workersCount := 5
	// maxErrorsCount здесь не влияет, так как ошибок нет.
	maxErrorsCount := 1

	start := time.Now()
	err := Run(tasks, workersCount, maxErrorsCount)
	elapsedTime := time.Since(start)
	require.NoError(t, err)
	// Читаем атомарно:
	require.Equal(t, int32(tasksCount), atomic.LoadInt32(&runTasksCount), "not all tasks were completed")
	// Если бы задачи выполнялись последовательно, время было бы ~sumSleepTime.
	// Проверяем, что оно существенно меньше.
	require.LessOrEqual(t, int64(elapsedTime), int64(sumSleepTime/2), "tasks were run sequentially?")
}

// TestRunMNonPositive проверяет сценарий m <= 0.
// В данной реализации m <= 0 трактуется как "игнорировать ошибки" – функция возвращает nil,
// а все задачи выполняются.
func TestRunMNonPositive(t *testing.T) {
	defer goleak.VerifyNone(t)

	tasksCount := 50
	tasks := make([]Task, 0, tasksCount)
	var runTasksCount int32

	for i := 0; i < tasksCount; i++ {
		tasks = append(tasks, func() error {
			atomic.AddInt32(&runTasksCount, 1)
			return errors.New("error")
		})
	}

	workersCount := 5
	err := Run(tasks, workersCount, 0)
	require.NoError(t, err, "no error expected when m <= 0")
	require.Equal(t, int32(tasksCount), atomic.LoadInt32(&runTasksCount), "all tasks should be executed when m <= 0")
}

// TestRunSuccessfulTasks проверяет сценарий, когда количество ошибок меньше порогового значения m.
// В этом случае все задачи должны выполниться, а функция вернуть nil.
func TestRunSuccessfulTasks(t *testing.T) {
	defer goleak.VerifyNone(t)

	tasksCount := 10
	tasks := make([]Task, 0, tasksCount)
	var runTasksCount int32

	// Пусть только одна задача возвращает ошибку, остальные – успешно.
	for i := 0; i < tasksCount; i++ {
		if i == 5 {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return errors.New("error")
			})
		} else {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}
	}

	workersCount := 3
	maxErrorsCount := 3
	err := Run(tasks, workersCount, maxErrorsCount)
	require.NoError(t, err)
	require.Equal(t, int32(tasksCount), atomic.LoadInt32(&runTasksCount), "not all tasks were executed")
}

// TestRunConcurrency проверяет, что задачи действительно выполняются параллельно,
// без использования time.Sleep для имитации задержки.
func TestRunConcurrency(t *testing.T) {
	defer goleak.VerifyNone(t)

	tasksCount := 20
	tasks := make([]Task, 0, tasksCount)
	var currentConcurrent int32
	var maxConcurrent int32
	blockCh := make(chan struct{})

	for i := 0; i < tasksCount; i++ {
		tasks = append(tasks, func() error {
			cur := atomic.AddInt32(&currentConcurrent, 1)
			// Обновляем максимум одновременных задач.
			for {
				mx := atomic.LoadInt32(&maxConcurrent)
				if cur > mx {
					if atomic.CompareAndSwapInt32(&maxConcurrent, mx, cur) {
						break
					}
				} else {
					break
				}
			}
			<-blockCh
			atomic.AddInt32(&currentConcurrent, -1)
			return nil
		})
	}

	workersCount := 5
	maxErrorsCount := 1

	doneCh := make(chan error, 1)
	go func() {
		err := Run(tasks, workersCount, maxErrorsCount)
		doneCh <- err
	}()

	// Проверяем, что в какой-то момент максимальное число одновременно выполняемых задач достигло количества воркеров.
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&maxConcurrent) >= int32(workersCount)
	}, time.Second, 10*time.Millisecond, "max concurrent tasks should be equal to workers count")

	close(blockCh)
	err := <-doneCh
	require.NoError(t, err)
}
