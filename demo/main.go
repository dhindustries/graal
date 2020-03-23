package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/dhindustries/graal/controller"
	"github.com/dhindustries/graal/prefab"
	"github.com/dhindustries/graal/utils"

	"github.com/dhindustries/graal/pathfinder"

	"github.com/dhindustries/graal/action"
	"github.com/dhindustries/graal/video"

	"github.com/dhindustries/graal/memory"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/dhindustries/graal/opengl"

	"github.com/dhindustries/graal/glfw"

	"github.com/dhindustries/graal/core"

	"github.com/dhindustries/graal"
)

type Application struct {
	graal.Application

	playerNode  graal.Node
	playerPool  action.AsyncPool
	playerQueue action.Queue

	camera       graal.Camera
	cameraPool   action.AsyncPool
	cameraQueue  action.Queue
	cameraFollow action.Switch
	keys         *controller.Keys

	tilemap        graal.Tilemap
	tilemapBuilder pathfinder.TilemapBuilder

	renderList  []interface{}
	disposeList []interface{}
}

func (app *Application) Prepare() error {
	app.disposeList = make([]interface{}, 0)
	app.renderList = make([]interface{}, 0)
	app.keys = controller.NewKeys(app.Keyboard())

	if err := app.initScene(); err != nil {
		return err
	}
	if err := app.initMap(); err != nil {
		return err
	}
	if err := app.initPlayer(); err != nil {
		return err
	}
	return nil
}

func (app *Application) Update(dt time.Duration) {
	if app.Keyboard().IsPressed(graal.KeyEscape) {
		app.Close()
	}
	if app.Mouse().IsPressed(graal.MouseButtonLeft) {
		if app.playerQueue.Empty() {
			x, y := app.mousePos()
			if x > 0 && y > 0 {
				x := uint(x)
				y := uint(y)
				app.moveToPath(x, y)
			}
		}
	}
	if app.Mouse().IsPressed(graal.MouseButtonRight) {
		x, y := app.mousePos()
		if x > 0 && y > 0 {
			x := uint(x)
			y := uint(y)
			app.tilemap.SetTile(x, y, (app.tilemap.Tile(x, y)+1)%4)
		}
	}
	app.keys.Update(dt)
	if app.playerQueue.Run(app.playerNode, dt) {
		app.playerPool.Run(app.playerNode, dt)
	}
	if app.cameraQueue.Run(app.camera, dt) {
		app.cameraPool.Run(app.camera, dt)
	}
}

func (app *Application) Render() {
	for _, r := range app.renderList {
		app.Renderer().Render(r)
	}
}

func (app *Application) Dispose() {
	for _, d := range app.renderList {
		graal.Release(d)
	}
	for _, d := range app.disposeList {
		memory.Dispose(d)
	}
}

func (app *Application) mousePos() (int, int) {
	m := app.Mouse().Cursor()
	iw := 1.0 / 16.0
	ih := 1.0 / 12.0
	cp := app.camera.Position()
	left, _, top, _ := app.camera.(graal.OrthoCamera).Viewport()
	x := int(m[0]/iw + cp[0] + left)
	y := int(m[1]/ih + cp[1] + top)
	return x, y
}

func (app *Application) initScene() error {
	prog, err := app.Resources().LoadProgram("assets/shader")
	if err != nil {
		return err
	}
	defer prog.Release()
	app.Renderer().BindProgram(prog)

	cam, err := app.Factory().OrthoCamera()
	if err != nil {
		return err
	}
	cam.SetNear(0)
	cam.SetFar(100)
	cam.SetPosition(mgl64.Vec3{0, 0, 15})
	cam.SetLookAt(mgl64.Vec3{0, 0, 0})
	cam.SetUp(mgl64.Vec3{0, 1, 0})
	cam.SetViewport(-8, -6, 8, 6)
	app.camera = cam
	app.disposeList = append(app.disposeList, cam)

	dpad := controller.NewDPad(app.Keyboard())
	dpad.SetSpeed(1)
	app.cameraPool.Add(&action.Update{Fn: dpad.Update})

	app.Renderer().SetCamera(cam)
	return nil
}

func (app *Application) initPlayer() error {
	t, err := app.loadTransparentTexture("assets/hero.png")
	if err != nil {
		return err
	}
	defer t.Release()
	q, err := app.Factory().Quad(0, 0, 1, 1)
	if err != nil {
		return err
	}
	defer q.Release()
	app.playerNode, err = app.Factory().Node()
	if err != nil {
		return nil
	}
	q.SetTexture(t)
	app.playerNode.Attach(q)
	app.playerNode.SetOrigin(mgl64.Vec3{-0.5, -0.5, 1})
	app.playerNode.SetPosition(mgl64.Vec3{0.5, 0.5, 1})
	app.renderList = append(app.renderList, app.playerNode)

	app.cameraFollow.Action = action.Func(func(target interface{}, dt time.Duration) bool {
		lookAt := app.playerNode.Position()
		pos := mgl64.Vec3{lookAt[0], lookAt[1], 15}
		utils.SetPosition(target, pos)
		utils.SetLookAt(target, lookAt)
		return false
	})
	app.cameraPool.Add(&app.cameraFollow)

	app.keys.Bind(graal.KeyO, func() {
		app.cameraFollow.Disable()
	})
	app.keys.Bind(graal.KeyP, func() {
		app.cameraFollow.Enable()
	})

	dpad := controller.NewDPad(app.Keyboard())
	dpad.SetSpeed(1)
	dpad.SetUpKey(graal.KeyW)
	dpad.SetDownKey(graal.KeyS)
	dpad.SetLeftKey(graal.KeyA)
	dpad.SetRightKey(graal.KeyD)
	app.playerPool.Add(&action.Update{Fn: dpad.Update})

	// app.s = action.NewScheduler()
	// app.disposeList = append(app.disposeList, app.s)
	// go app.s.Start(app.n, 60)
	return nil
}

func (app *Application) moveToPath(x, y uint) {
	pos := app.playerNode.Position()
	anim := action.Sequence{}
	path := app.tilemapBuilder.Build(app.tilemap).FindPath(int(pos[0]), int(pos[1]), int(x), int(y))
	for _, coords := range path {
		anim.Add(&action.MoveTo{
			Position: mgl64.Vec3{float64(coords[0]) + 0.5, float64(coords[1]) + 0.5, pos[2]},
			Duration: time.Second,
		})
	}
	app.playerQueue.Add(&anim)
}

func (app *Application) initMap() error {
	app.tilemapBuilder.Tiles = map[uint]pathfinder.TileInfo{}
	m, err := app.loadMap("assets/map.xml")
	if err != nil {
		return err
	}
	app.tilemap = m
	app.renderList = append(app.renderList, m)
	return nil
}

func (app *Application) loadMap(m string) (graal.Tilemap, error) {
	tilemapPrefab, err := app.Resources().LoadPrefab(m)
	if err != nil {
		return nil, err
	}
	defer tilemapPrefab.Release()
	tilemapHandle, err := tilemapPrefab.Spawn()
	if err != nil {
		return nil, err
	}
	return tilemapHandle.(graal.Tilemap), nil
}

func (app *Application) loadTileset(m string) (graal.Tileset, error) {
	tex, err := app.loadTransparentTexture(m)
	if err != nil {
		return nil, err
	}
	defer tex.Release()
	ts, err := app.Factory().Tileset()
	if err != nil {
		return nil, err
	}
	ts.SetTexture(tex)
	ts.SetTileSize(16, 16)
	return ts, nil
}

func (app *Application) loadTransparentTexture(m string) (graal.Texture, error) {
	img, err := app.Resources().LoadImage(m)
	if err != nil {
		return nil, err
	}
	defer img.Release()
	imgKeyFilter(img, graal.Color{1, 0, 1, 0})
	tex, err := app.Factory().Texture(img.Size())
	tex.SetMode(graal.TextureModeSharp)
	if err != nil {
		return nil, err
	}
	if err := tex.Draw(img); err != nil {
		tex.Release()
		return nil, err
	}
	return tex, nil
}

func imgKeyFilter(img graal.Image, k graal.Color) {
	eqf := func(a, b float64) bool {
		d := math.Abs(b - a)
		return d <= 0.000001
	}
	eqc := func(a, b graal.Color) bool {
		return eqf(a[0], b[0]) && eqf(a[1], b[1]) && eqf(a[2], b[2])
	}
	w, h := img.Size()
	for y := uint(0); y < h; y++ {
		for x := uint(0); x < w; x++ {
			if eqc(img.Get(x, y), k) {
				img.Set(x, y, graal.Color{0, 0, 0, 0})
			}
		}
	}
}

func displayInfo() chan bool {
	mb := func(v uint64) string {
		return fmt.Sprintf("%v MB", v/1024/1024)
	}
	print := func() {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		fmt.Printf(
			"<go routines: %v alloc: %v heap: %v sys: %v total: %s>\n",
			runtime.NumGoroutine()-3,
			mb(mem.Alloc),
			mb(mem.HeapAlloc),
			mb(mem.Sys),
			mb(mem.TotalAlloc),
		)
	}
	brk := make(chan bool)
	end := make(chan bool, 1)
	clk := time.NewTicker(5 * time.Millisecond)
	go func() {
		for range brk {
		}
		end <- true
	}()
	go func() {
	loop:
		for {
			print()
			select {
			case <-clk.C:
				break
			case <-brk:
				break loop
			}
		}
	}()
	return brk
}

func main() {
	// defer close(displayInfo())
	err := graal.Bootstrap(
		&core.Provider{},
		&glfw.Provider{},
		&video.Provider{},
		&opengl.Provider{},
		&prefab.Provider{},
	).Run(&Application{})
	if err != nil {
		panic(err)
	}
}
