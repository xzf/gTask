package task

func New(limit int) *task {
	task := &task{
		limit:          limit,
		jobQueue:       nil,
		jobControlChan: make(chan uint8, limit),
		status:         0,
		waitChan:       make(chan uint8, 1),
	}
	go task.asyncRunThread()
	return task
}
