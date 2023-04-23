package task

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	n := int64(0)
	task := New(1)
	times := int64(10000)
	for i := int64(0); i < times; i++ {
		task.Run(func() {
			time.Sleep(time.Microsecond * 10)
			atomic.AddInt64(&n, 1)
		})
	}
	task.WaitNotJob()
	if times != n {
		panic("85ho9885nj")
	}
}
