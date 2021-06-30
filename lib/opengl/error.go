package opengl

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)

func logGlMessage(
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
	fmt.Printf("OpenGL %s %s from %s: %s\n", types[gltype], serverities[severity], sources[source], message)
}

func glError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv, getObjInfoLogFn getObjInfoLog, failMsg string) error {

	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)

		log := gl.Str(strings.Repeat("\x00", int(logLength)+1))
		getObjInfoLogFn(glHandle, logLength, nil, log)
		return fmt.Errorf("%s: %s", failMsg, gl.GoStr(log))
	}

	return nil
}
