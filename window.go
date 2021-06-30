package graal

type apiWindow interface {
	WindowSize() (width, height uint)
	initWindow() error
	finishWindow()
	openWindow(width, height uint, title string) error
	closeWindow()
	isWindowOpen() bool
	updateWindow()
	renderWindow()
}

type protoWindow struct {
	WindowSize   func(api Api) (width, height uint)
	InitWindow   func(api Api) error
	FinishWindow func(api Api)
	OpenWindow   func(api Api, width, height uint, title string) error
	CloseWindow  func(api Api)
	IsWindowOpen func(api Api) bool
	UpdateWindow func(api Api)
	RenderWindow func(api Api)
}

func WindowSize() (width, height uint) {
	return api.WindowSize()
}

func (api *apiAdapter) WindowSize() (width, height uint) {
	if api.proto.WindowSize == nil {
		panic("api.WindowSize is not implemented")
	}
	return api.proto.WindowSize(api)
}

func (api *apiAdapter) initWindow() error {
	//if api.proto.InitWindow == nil {
	//	panic("api.InitWindow is not implemented")
	//}
	//return api.proto.InitWindow(api)
	if api.proto.InitWindow != nil {
		return api.proto.InitWindow(api)
	}
	return nil
}

func (api *apiAdapter) finishWindow() {
	//if api.proto.FinishWindow == nil {
	//	panic("api.FinishWindow is not implemented")
	//}
	//api.proto.FinishWindow(api)
	if api.proto.FinishWindow != nil {
		api.proto.FinishWindow(api)
	}
}

func (api *apiAdapter) openWindow(width, height uint, title string) error {
	if api.proto.OpenWindow == nil {
		panic("api.OpenWindow is not implemented")
	}
	return api.proto.OpenWindow(api, width, height, title)
}

func Close() {
	api.closeWindow()
}

func (api *apiAdapter) closeWindow() {
	if api.proto.CloseWindow == nil {
		panic("api.CloseWindow is not implemented")
	}
	api.proto.CloseWindow(api)
}

func (api *apiAdapter) isWindowOpen() bool {
	if api.proto.IsWindowOpen == nil {
		panic("api.IsWindowOpen is not implemented")
	}
	return api.proto.IsWindowOpen(api)
}

func (api *apiAdapter) updateWindow() {
	if api.proto.UpdateWindow == nil {
		panic("api.UpdateWindow is not implemented")
	}
	api.proto.UpdateWindow(api)
}

func (api *apiAdapter) renderWindow() {
	if api.proto.RenderWindow == nil {
		panic("api.RenderWindow is not implemented")
	}
	api.proto.RenderWindow(api)
}
