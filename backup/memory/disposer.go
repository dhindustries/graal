package memory

type Disposer interface {
	Dispose()
}

func Dispose(v interface{}) {
	if x, ok := v.(Disposer); ok {
		x.Dispose()
	}
}
