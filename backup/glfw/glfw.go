package glfw

import (
	"sync"

	"github.com/go-gl/glfw/v3.2/glfw"
)

var windowCounter uint = 0
var windowCounterLock sync.Mutex

func glfwInit() {
	windowCounterLock.Lock()
	defer windowCounterLock.Unlock()
	if windowCounter == 0 {
		glfw.Init()
	}
	windowCounter++
}

func glfwFinit() {
	windowCounterLock.Lock()
	defer windowCounterLock.Unlock()
	if windowCounter == 1 {
		glfw.Terminate()
	}
	windowCounter--
}
