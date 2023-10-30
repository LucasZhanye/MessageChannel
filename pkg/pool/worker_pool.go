package pool

import (
	"sync"
	"time"

	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"go.uber.org/atomic"
)

const (
	taskIdleTimeout = 5 * time.Second
)

type Task struct {
	Fn    func(...any)
	Param []any
}

type WorkerPool struct {
	maxWorker       uint32
	currentWorker   atomic.Uint32
	tasksQueue      chan *Task
	workersQueue    chan *Task
	waitingQueue    *linkedlistqueue.Queue
	workerWaitGroup sync.WaitGroup
	stopChan        chan struct{}
	closeChan       chan struct{}
	closed          bool
	closeLock       sync.Mutex
}

func NewWorkerPool(maxWorker uint32) *WorkerPool {
	if maxWorker <= 0 {
		maxWorker = 1
	}

	return &WorkerPool{
		maxWorker:    maxWorker,
		tasksQueue:   make(chan *Task),
		workersQueue: make(chan *Task),
		waitingQueue: linkedlistqueue.New(),
		stopChan:     make(chan struct{}),
		closeChan:    make(chan struct{}),
	}
}

func (wp *WorkerPool) Put(task *Task) {
	if task.Fn != nil {
		wp.tasksQueue <- task
	}
}

func (wp *WorkerPool) GetWaitingTaskSize() int {
	return wp.waitingQueue.Size()
}

func (wp *WorkerPool) GetCurrentWorkerSize() uint32 {
	return wp.currentWorker.Load()
}

func (wp *WorkerPool) Run() {
	t := time.NewTimer(taskIdleTimeout)

LOOP:
	for {
		// check waiting queue
		if wp.waitingQueue.Size() > 0 {
			if !wp.handleWaitingTask() {
				break LOOP
			}
			// if waiting queue not null, handle waiting task at first
			continue
		}

		select {
		case task := <-wp.tasksQueue:

			// send task to workers
			select {
			case wp.workersQueue <- task:
			default:
				// all workers are busy, so need to start a new worker
				if wp.currentWorker.Load() < wp.maxWorker {
					go wp.worker(task)

					wp.currentWorker.Inc()
				} else {
					// out of maxWorker, save to waiting queue
					wp.waitingQueue.Enqueue(task)
				}
			}
		case <-t.C:
			if wp.currentWorker.Load() > 0 {
				if wp.stopOneWorker() {
					wp.currentWorker.Dec()
				}
			}
			t.Reset(taskIdleTimeout)
		case <-wp.closeChan:
			break LOOP
		}

	}

	// waiting all task finish
	wp.execAllWaitingTask()

	// if worker pool closed, stop all worker
	wp.stopAllWorker()

	// wait all worker goroutinue exit
	wp.workerWaitGroup.Wait()

	t.Stop()

	// notify stop finish
	close(wp.stopChan)
}

func (wp *WorkerPool) worker(task *Task) {
	targetTask := task

	wp.workerWaitGroup.Add(1)

	for {
		// exec task
		targetTask.Fn(targetTask.Param...)

		// wait task
		targetTask = <-wp.workersQueue

		if targetTask == nil {
			wp.workerWaitGroup.Done()
			return
		}
	}
}

func (wp *WorkerPool) stopOneWorker() bool {
	select {
	// send nil to worker and then worker will close
	case wp.workersQueue <- nil:
		return true
	default:
		return false
	}
}

func (wp *WorkerPool) stopAllWorker() {
	for wp.currentWorker.Load() > 0 {
		wp.workersQueue <- nil
		wp.currentWorker.Dec()
	}
}

func (wp *WorkerPool) handleWaitingTask() bool {
	// get task from waiting queue, send it to the worker
	waitTask, _ := wp.waitingQueue.Peek()
	wtask := waitTask.(*Task)

	select {
	case task := <-wp.tasksQueue:
		// if have a new task ,add to waiting quene
		wp.waitingQueue.Enqueue(task)
	case wp.workersQueue <- wtask:
		// when worker receive success, delete task from waiting queue
		wp.waitingQueue.Dequeue()
	case <-wp.closeChan:
		return false
	}

	return true
}

func (wp *WorkerPool) execAllWaitingTask() {
	for wp.waitingQueue.Size() > 0 {

		waitTask, _ := wp.waitingQueue.Peek()
		task := (waitTask).(*Task)

		wp.workersQueue <- task

		wp.waitingQueue.Dequeue()
	}
}

func (wp *WorkerPool) Stop() {
	wp.closeLock.Lock()

	if wp.closed {
		wp.closeLock.Unlock()
		return
	}
	wp.closed = true
	close(wp.closeChan)

	wp.closeLock.Unlock()

	<-wp.stopChan
}
