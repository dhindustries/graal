package queue

type Invoker struct {
	Q *Queue
}

func NewInvoker(q *Queue) *Invoker {
	return &Invoker{Q: q}
}

func (invoker Invoker) Call(fn func()) {
	done := make(chan interface{}, 1)
	invoker.Q.Push(func() {
		fn()
		done <- nil
	})
	<-done
	close(done)
}
