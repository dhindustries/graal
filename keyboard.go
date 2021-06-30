package graal

//const _Ciconst_GLFW_LOSE_CONTEXT_ON_RESET = 0x31002
//const _Ciconst_GLFW_MAXIMIZED = 0x20008
//const _Ciconst_GLFW_MOD_ALT = 0x4
//const _Ciconst_GLFW_MOD_CONTROL = 0x2
//const _Ciconst_GLFW_MOD_SHIFT = 0x1
//const _Ciconst_GLFW_MOD_SUPER = 0x8
//const _Ciconst_GLFW_MOUSE_BUTTON_1 = 0x0
//const _Ciconst_GLFW_MOUSE_BUTTON_2 = 0x1
//const _Ciconst_GLFW_MOUSE_BUTTON_3 = 0x2
//const _Ciconst_GLFW_MOUSE_BUTTON_4 = 0x3
//const _Ciconst_GLFW_MOUSE_BUTTON_5 = 0x4
//const _Ciconst_GLFW_MOUSE_BUTTON_6 = 0x5
//const _Ciconst_GLFW_MOUSE_BUTTON_7 = 0x6
//const _Ciconst_GLFW_MOUSE_BUTTON_8 = 0x7
//const _Ciconst_GLFW_MOUSE_BUTTON_LAST = 0x7
//const _Ciconst_GLFW_MOUSE_BUTTON_LEFT = 0x0
//const _Ciconst_GLFW_MOUSE_BUTTON_MIDDLE = 0x2
//const _Ciconst_GLFW_MOUSE_BUTTON_RIGHT = 0x1

type Key int

const (
	KeyUnknown      Key = -0x1
	KeySpace        Key = 0x20
	KeyApostrophe   Key = 0x27
	KeyComma        Key = 0x2c
	KeyMinus        Key = 0x2d
	KeyPeriod       Key = 0x2e
	KeySlash        Key = 0x2f
	Key0            Key = 0x30
	Key1            Key = 0x31
	Key2            Key = 0x32
	Key3            Key = 0x33
	Key4            Key = 0x34
	Key5            Key = 0x35
	Key6            Key = 0x36
	Key7            Key = 0x37
	Key8            Key = 0x38
	Key9            Key = 0x39
	KeySemicolon    Key = 0x3b
	KeyEqual        Key = 0x3d
	KeyA            Key = 0x41
	KeyB            Key = 0x42
	KeyC            Key = 0x43
	KeyD            Key = 0x44
	KeyE            Key = 0x45
	KeyF            Key = 0x46
	KeyG            Key = 0x47
	KeyH            Key = 0x48
	KeyI            Key = 0x49
	KeyJ            Key = 0x4a
	KeyK            Key = 0x4b
	KeyL            Key = 0x4c
	KeyM            Key = 0x4d
	KeyN            Key = 0x4e
	KeyO            Key = 0x4f
	KeyP            Key = 0x50
	KeyQ            Key = 0x51
	KeyR            Key = 0x52
	KeyS            Key = 0x53
	KeyT            Key = 0x54
	KeyU            Key = 0x55
	KeyV            Key = 0x56
	KeyW            Key = 0x57
	KeyX            Key = 0x58
	KeyY            Key = 0x59
	KeyZ            Key = 0x5a
	KeyLeftBracket  Key = 0x5b
	KeyBackslash    Key = 0x5c
	KeyRightBracket Key = 0x5d
	KeyGraveAccent  Key = 0x60
	KeyWorld1       Key = 0xa1
	KeyWorld2       Key = 0xa2
	KeyEscape       Key = 0x100
	KeyEnter        Key = 0x101
	KeyTab          Key = 0x102
	KeyBackspace    Key = 0x103
	KeyInsert       Key = 0x104
	KeyDelete       Key = 0x105
	KeyRight        Key = 0x106
	KeyLeft         Key = 0x107
	KeyDown         Key = 0x108
	KeyUp           Key = 0x109
	KeyPageUp       Key = 0x10a
	KeyPageDown     Key = 0x10b
	KeyHome         Key = 0x10c
	KeyEnd          Key = 0x10d
	KeyCapsLock     Key = 0x118
	KeyScrollLock   Key = 0x119
	KeyNumLock      Key = 0x11a
	KeyPrintScreen  Key = 0x11b
	KeyPause        Key = 0x11c
	KeyF1           Key = 0x122
	KeyF2           Key = 0x123
	KeyF3           Key = 0x124
	KeyF4           Key = 0x125
	KeyF5           Key = 0x126
	KeyF6           Key = 0x127
	KeyF7           Key = 0x128
	KeyF8           Key = 0x129
	KeyF9           Key = 0x12a
	KeyF10          Key = 0x12b
	KeyF11          Key = 0x12c
	KeyF12          Key = 0x12d
	KeyF13          Key = 0x12e
	KeyF14          Key = 0x12f
	KeyF15          Key = 0x130
	KeyF16          Key = 0x131
	KeyF17          Key = 0x132
	KeyF18          Key = 0x133
	KeyF19          Key = 0x134
	KeyF20          Key = 0x135
	KeyF21          Key = 0x136
	KeyF22          Key = 0x137
	KeyF23          Key = 0x138
	KeyF24          Key = 0x139
	KeyF25          Key = 0x13a
	KeyKP0          Key = 0x140
	KeyKP1          Key = 0x141
	KeyKP2          Key = 0x142
	KeyKP3          Key = 0x143
	KeyKP4          Key = 0x144
	KeyKP5          Key = 0x145
	KeyKP6          Key = 0x146
	KeyKP7          Key = 0x147
	KeyKP8          Key = 0x148
	KeyKP9          Key = 0x149
	KeyKPDecimal    Key = 0x14a
	KeyKPDivide     Key = 0x14b
	KeyKPMultiply   Key = 0x14c
	KeyKPSubtract   Key = 0x14d
	KeyKPAdd        Key = 0x14e
	KeyKPEnter      Key = 0x14f
	KeyKPEqual      Key = 0x150
	KeyLeftShift    Key = 0x154
	KeyLeftControl  Key = 0x155
	KeyLeftAlt      Key = 0x156
	KeyLeftSuper    Key = 0x157
	KeyRightShift   Key = 0x158
	KeyRightControl Key = 0x159
	KeyRightAlt     Key = 0x15a
	KeyRightSuper   Key = 0x15b
	KeyMenu         Key = 0x15c
	KeyLast         Key = 0x15c
)

type apiKeyboard interface {
	IsKeyDown(key Key) bool
	IsKeyUp(key Key) bool
	IsKeyPressed(key Key) bool
	IsKeyReleased(key Key) bool
	KeyboardInput() (data <-chan rune, close func())
}

type protoKeyboard struct {
	IsKeyDown func(api Api, key Key) bool
	IsKeyUp func(api Api, key Key) bool
	IsKeyPressed func(api Api, key Key) bool
	IsKeyReleased func(api Api, key Key) bool
	KeyboardInput func(api Api) (data <-chan rune, close func())
}

func IsKeyDown(key Key) bool {
	return api.IsKeyDown(key)
}

func (api *apiAdapter) IsKeyDown(key Key) bool {
	if api.proto.IsKeyDown == nil {
		panic("api.IsKeyDown is not implemented")
	}
	return api.proto.IsKeyDown(api, key)
}

func IsKeyUp(key Key) bool {
	return api.IsKeyUp(key)
}

func (api *apiAdapter) IsKeyUp(key Key) bool {
	if api.proto.IsKeyUp == nil {
		panic("api.IsKeyUp is not implemented")
	}
	return api.proto.IsKeyUp(api, key)
}

func IsKeyPressed(key Key) bool {
	return api.IsKeyPressed(key)
}

func (api *apiAdapter) IsKeyPressed(key Key) bool {
	if api.proto.IsKeyPressed == nil {
		panic("api.IsKeyPressed is not implemented")
	}
	return api.proto.IsKeyPressed(api, key)
}

func IsKeyReleased(key Key) bool {
	return api.IsKeyReleased(key)
}

func (api *apiAdapter) IsKeyReleased(key Key) bool {
	if api.proto.IsKeyReleased == nil {
		panic("api.IsKeyReleased is not implemented")
	}
	return api.proto.IsKeyReleased(api, key)
}

func KeyboardInput() (data <-chan rune, close func()) {
	return api.KeyboardInput()
}

func (api *apiAdapter) KeyboardInput() (data <-chan rune, close func()) {
	if api.proto.KeyboardInput == nil {
		panic("api.KeyboardInput is not implemented")
	}
	return api.proto.KeyboardInput(api)
}
