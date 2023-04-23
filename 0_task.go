package task

import (
	"sync"
	"sync/atomic"
)

const (
	statusWait int32 = 0
	statusRun  int32 = 1
)

type Task struct {
	limit          int
	jobLock        sync.Mutex
	jobQueue       []func()
	jobControlChan chan uint8
	jobWg          sync.WaitGroup
	status         int32
	waitChan       chan uint8
}

//type job struct {
//}

func (t *Task) Run(job func()) {
	if job == nil {
		return
	}
	t.jobLock.Lock()
	t.jobQueue = append(t.jobQueue, job)
	t.jobWg.Add(1)
	if atomic.CompareAndSwapInt32(&t.status, statusWait, statusRun) {
		t.waitChan <- 0
	}
	t.jobLock.Unlock()
}

func (t *Task) asyncRunThread() {
	for {
		<-t.waitChan
		t.waitChan <- 0
		t.jobLock.Lock()
		var job func()
		if len(t.jobQueue) == 0 {
			atomic.CompareAndSwapInt32(&t.status, statusRun, statusWait)
			<-t.waitChan
			t.jobLock.Unlock()
			continue
		}
		job = t.jobQueue[0]
		t.jobQueue = t.jobQueue[1:]
		t.jobLock.Unlock()

		t.jobControlChan <- 1
		go func() {
			defer func() {
				t.jobWg.Done()
				<-t.jobControlChan
			}()
			job()
		}()
	}
}

func (t *Task) WaitNotJob() {
	t.jobWg.Wait()
}
