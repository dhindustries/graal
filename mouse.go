package graal

type apiMouse interface {
	MousePosition() (x, y uint)
}

type protoMouse struct {
	MousePosition func(api Api) (x, y uint)
}

func MousePosition() (x, y uint) {
	return api.MousePosition()
}

func (api *apiAdapter) MousePosition() (x, y uint) {
	if api.proto.MousePosition == nil {
		panic("api.MousePosition is not implemented")
	}
	return api.proto.MousePosition(api)
}
