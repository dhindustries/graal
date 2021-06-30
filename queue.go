package graal

type Queue struct {
	tasks chan func()
}

func (queue *Queue) Start() {
	if queue.tasks != nil {
		panic("queue is already started")
	}
	queue.tasks = make(chan func(), 8)
}

func (queue *Queue) Stop() {
	if queue.tasks != nil {
		close(queue.tasks)
		queue.tasks = nil
	}
}

func (queue *Queue) Pull() {
	if queue.tasks == nil {
		panic("queue is not started")
	}
	for task := range queue.tasks {
		task()
	}
}

func (queue *Queue) Enqueue(fn func()) {
	if queue.tasks == nil {
		panic("queue is not started")
	}
	queue.tasks <- fn
}
