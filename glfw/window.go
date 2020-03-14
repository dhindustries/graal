package glfw

import (
	"fmt"
	"sync"

	"github.com/go-gl/glfw/v3.2/glfw"
)

var windowCounter uint = 0
var windowCounterLock sync.Mutex

type Window struct {
	Title         string
	Width, Height uint
	Handle        *glfw.Window
	lock          sync.Mutex
}

func glfwInitialize() {
	windowCounterLock.Lock()
	defer windowCounterLock.Unlock()
	if windowCounter == 0 {
		glfw.Init()
	}
	windowCounter++
}

func glfwTerminate() {
	windowCounterLock.Lock()
	defer windowCounterLock.Unlock()
	if windowCounter == 1 {
		glfw.Terminate()
	}
	windowCounter--
}

func NewWindow(width, height uint, title string) *Window {
	w := &Window{}
	w.Width = width
	w.Height = height
	w.Title = title
	return w
}

func (window *Window) Open() error {
	var err error
	window.lock.Lock()
	defer window.lock.Unlock()
	if window.Handle != nil {
		return fmt.Errorf("Window is already open")
	}
	glfwInitialize()
	window.Handle, err = glfw.CreateWindow(int(window.Width), int(window.Height), window.Title, nil, nil)
	if err != nil {
		glfwTerminate()
		return err
	}
	window.Handle.MakeContextCurrent()
	window.Handle.Show()

	return nil
}

func (window *Window) Close() {
	window.lock.Lock()
	defer window.lock.Unlock()
	if window.Handle != nil {
		window.Handle.SetShouldClose(true)
	}
}

func (window *Window) Dispose() {
	window.lock.Lock()
	defer window.lock.Unlock()
	if window.Handle == nil {
		panic(fmt.Errorf("Cannot dispose window: window is not created"))
	}
	window.Handle.Destroy()
	glfwTerminate()
}

func (window *Window) IsOpen() bool {
	window.lock.Lock()
	defer window.lock.Unlock()

	return window.Handle != nil && !window.Handle.ShouldClose()
}

func (window *Window) PullMessages() {
	window.lock.Lock()
	defer window.lock.Unlock()

	if window.Handle != nil {
		glfw.PollEvents()
	}
}

func (window *Window) Swap() {
	window.lock.Lock()
	defer window.lock.Unlock()

	if window.Handle != nil {
		window.Handle.SwapBuffers()
	}
}
