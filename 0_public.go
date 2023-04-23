package task

func NewTask(limit int) *Task {
	task := &Task{
		limit:          limit,
		jobQueue:       nil,
		jobControlChan: make(chan uint8, limit),
		status:         0,
	}
	go task.asyncRunThread()
	return task
}
