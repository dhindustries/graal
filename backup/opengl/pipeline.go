package opengl

// import (
// 	"github.com/dhindustries/graal"
// )

// type pipeline struct {
// 	camera  graal.Camera
// 	program graal.Program
// }

// func (pipeline *pipeline) SetCamera(camera graal.Camera) {
// 	pipeline.camera = camera
// }
// func (pipeline *pipeline) Camera() graal.Camera {
// 	return pipeline.camera

// }
// func (pipeline *pipeline) SetProgram(program graal.Program) {
// 	if pipeline.program != nil {
// 		pipeline.program.Release()
// 	}
// 	if program != nil {
// 		program.Acquire()
// 	}
// 	pipeline.program = program
// }

// func (pipeline *pipeline) Program() graal.Program {
// 	return pipeline.program
// }

// func (pipeline *pipeline) Dispose() {
// 	pipeline.SetProgram(nil)
// }

// func (renderer *renderer) applyPipeline(pipeline graal.Pipeline) {
// 	cam := pipeline.Camera()
// 	prog := pipeline.Program()
// 	if prog != nil {
// 		if cam != nil {
// 			prog.SetView(cam.View())
// 			prog.SetProjection(cam.Projection())
// 		}
// 		renderer.applyProgram(prog)
// 	}
// }
