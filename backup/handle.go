package graal

type Handle interface {
	Acquire()
	Release()
	IsValid() bool
}

type acquirer interface {
	Acquire()
}

type releaser interface {
	Release()
}

func Acquire(v interface{}) {
	if v, ok := v.(acquirer); ok {
		v.Acquire()
	}
}

func Release(v interface{}) {
	if v, ok := v.(releaser); ok {
		v.Release()
	}
}
