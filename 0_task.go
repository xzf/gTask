package task

import (
	"sync"
)

const (
	statusWait = 0
	statusRun  = 1
)

type Task struct {
	limit          int
	jobLock        sync.Mutex
	jobQueue       []func()
	jobControlChan chan uint8
	status         int
}

//type job struct {
//}

func (t *Task) Run(job func()) {
	if job == nil {
		return
	}
	t.jobLock.Lock()

	t.jobQueue = append(t.jobQueue, job)
	if t.status == statusWait {
		t.status = statusRun
		t.jobControlChan <- 0
	}
	t.jobLock.Unlock()
}

func (t *Task) asyncRunThread() {
	for {
		<-t.jobControlChan
		t.jobLock.Lock()
		var job func()
		if len(t.jobQueue) != 0 {
			job = t.jobQueue[0]
			t.jobQueue = t.jobQueue[1:]
		}
		if job == nil {
			t.status = statusWait
			t.jobLock.Unlock()
			continue
		}
		t.status = statusRun
		t.jobLock.Unlock()

		go func() {
			defer func() {
				t.jobControlChan <- 0
			}()
			job()
		}()
	}
}
