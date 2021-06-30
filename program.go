package graal

type Program interface {
}

type apiProgram interface {
	NewProgram(vertex VertexShader, fragment FragmentShader) (Program, error)
	UseProgram(program Program)
}

type protoProgram struct {
	NewProgram func(api Api, vertex VertexShader, fragment FragmentShader) (Program, error)
	UseProgram func(api Api, program Program)
}

func NewProgram(vertex VertexShader, fragment FragmentShader) (Program, error) {
	return api.NewProgram(vertex, fragment)
}

func (api *apiAdapter) NewProgram(vertex VertexShader, fragment FragmentShader) (Program, error) {
	if api.proto.NewProgram == nil {
		panic("api.NewProgram is not implemented")
	}
	return api.proto.NewProgram(api, vertex, fragment)
}

func UseProgram(program Program) {
	api.UseProgram(program)
}

func (api *apiAdapter) UseProgram(program Program) {
	if api.proto.UseProgram == nil {
		panic("api.UseProgram is not implemented")
	}
	api.proto.UseProgram(api, program)
}
