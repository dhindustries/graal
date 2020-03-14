package graal

type Window interface {
	Open() error
	Close()
	Dispose()
	IsOpen() bool
	PullMessages()
}
