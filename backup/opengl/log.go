package opengl

import (
	"os"
	"strings"
	"unsafe"

	"github.com/dhindustries/graal"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type logger struct {
	api *graal.Api
}

func (l *logger) enabled() bool {
	return strings.ToLower(os.Getenv("GLLOG")) == "on"
}

func (l *logger) watch() {
	l.api.Invoke(func() {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
		gl.DebugMessageCallback(l.glHandle, nil)
		gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nil, true)
	})
}

func (l *logger) glHandle(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer,
) {
	sources := map[uint32]string{
		gl.DEBUG_SOURCE_API:             "API",
		gl.DEBUG_SOURCE_WINDOW_SYSTEM:   "Window System",
		gl.DEBUG_SOURCE_SHADER_COMPILER: "Shader Compiler",
		gl.DEBUG_SOURCE_THIRD_PARTY:     "Third Party",
		gl.DEBUG_SOURCE_APPLICATION:     "Application",
		gl.DEBUG_SOURCE_OTHER:           "Other",
	}
	types := map[uint32]string{
		gl.DEBUG_TYPE_ERROR:               "Error",
		gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR: "Deprecated Behaviour",
		gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR:  "Undefined Behaviour",
		gl.DEBUG_TYPE_PORTABILITY:         "Portability",
		gl.DEBUG_TYPE_PERFORMANCE:         "Performance",
		gl.DEBUG_TYPE_MARKER:              "Marker",
		gl.DEBUG_TYPE_PUSH_GROUP:          "Push Group",
		gl.DEBUG_TYPE_POP_GROUP:           "Pop Group",
		gl.DEBUG_TYPE_OTHER:               "Other",
	}
	serverities := map[uint32]string{
		gl.DEBUG_SEVERITY_HIGH:         "high",
		gl.DEBUG_SEVERITY_MEDIUM:       "medium",
		gl.DEBUG_SEVERITY_LOW:          "low",
		gl.DEBUG_SEVERITY_NOTIFICATION: "notification",
	}
	l.log("OpenGL %s %s from %s: %s\n", types[gltype], serverities[severity], sources[source], message)
}

func (l *logger) log(f string, v ...interface{}) {
	l.api.Logf(l.api, f, v...)
}
