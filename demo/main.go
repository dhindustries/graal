package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/dhindustries/graal/action"
	"github.com/dhindustries/graal/controller"
	"github.com/dhindustries/graal/video"

	"github.com/dhindustries/graal/memory"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/dhindustries/graal/opengl"

	"github.com/dhindustries/graal/glfw"

	"github.com/dhindustries/graal/core"

	"github.com/dhindustries/graal"
)

type Application struct {
	graal.Application
	cam         graal.OrthoCamera
	m           graal.Tilemap
	s           *action.Scheduler
	n           graal.Node
	renderList  []interface{}
	disposeList []interface{}
}

func (app *Application) Prepare() error {
	app.disposeList = make([]interface{}, 0)
	app.renderList = make([]interface{}, 0)

	if err := app.initScene(); err != nil {
		return err
	}
	if err := app.initMap(); err != nil {
		return err
	}
	if err := app.initPlayer(); err != nil {
		return err
	}

	dpad := controller.NewDPad(app.Keyboard())
	dpad.SetTarget(app.cam)
	dpad.SetSpeed(10)
	go dpad.Start(60.0)
	app.disposeList = append(app.disposeList, dpad)

	return nil
}

func (app *Application) Update(dt float32) {
	if app.Keyboard().IsPressed(graal.KeyEscape) {
		app.Close()
	}
	if app.Mouse().IsPressed(graal.MouseButtonLeft) {
		x, y := app.mouseTile()
		go app.moveToPath(mgl32.Vec3{float32(x), float32(y), 1})
	}
	if app.Mouse().IsPressed(graal.MouseButtonRight) {
		x, y := app.mouseTile()
		app.m.SetTile(x, y, (app.m.Tile(x, y)+1)%4)

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

func (app *Application) mouseTile() (uint, uint) {
	m := app.Mouse().Cursor()
	iw := float32(1.0 / 16.0)
	ih := float32(1.0 / 12.0)
	cp := app.cam.Position()
	x := uint(m[0]/iw + cp[0])
	y := uint(m[1]/ih + cp[1])
	return x, y
}

func (app *Application) initScene() error {
	prog, err := app.Resources().LoadProgram("assets/shader")
	if err != nil {
		return err
	}
	defer prog.Release()
	app.Renderer().BindProgram(prog)

	app.cam, err = app.Factory().OrthoCamera()
	if err != nil {
		return err
	}
	app.cam.SetNear(0)
	app.cam.SetFar(100)
	app.cam.SetPosition(mgl32.Vec3{0, 0, 15})
	app.cam.SetLookAt(mgl32.Vec3{0, 0, 0})
	app.cam.SetUp(mgl32.Vec3{0, 1, 0})
	app.cam.SetViewport(0, 0, 16, 12)
	app.disposeList = append(app.disposeList, app.cam)

	app.Renderer().SetCamera(app.cam)
	return nil
}

type nodeElem struct {
	w    float32
	x, y uint
	p    *nodeElem
}

func (app *Application) findPath(sx, sy, dx, dy uint) []mgl32.Vec3 {
	w, h := app.m.Size()
	nodes := make([]*nodeElem, w*h)
	q := make([]nodeElem, 0)
	visited := func(x, y int) bool {
		return nodes[uint(y)*w+uint(x)] != nil
	}
	available := func(x, y int) bool {
		t := app.m.Tile(uint(x), uint(y))
		return t == 1 || t == 3
	}
	dirs := [][]int{
		[]int{0, -1},
		[]int{-1, 0},
		[]int{0, 1},
		[]int{1, 0},
	}

	q = append(q, nodeElem{
		w: 0,
		x: sx, y: sy,
		p: nil,
	})
	var result *nodeElem

	for len(q) > 0 {
		node := q[0]
		if node.x == dx && node.y == dy {
			result = &node
			break
		}
		q = q[1:]
		nodes[node.y*w+node.x] = &node

		for _, dir := range dirs {
			nx, ny := int(node.x)+dir[0], int(node.y)+dir[1]
			if (nx < int(w) && ny < int(h) && nx >= 0 && ny >= 0) && (!visited(nx, ny) || available(nx, ny)) {
				q = append(q, nodeElem{
					w: node.w + 1,
					x: uint(nx), y: uint(ny),
					p: &node,
				})
			}
		}
	}
	steps := make([]mgl32.Vec3, 0)
	for result != nil {
		i := []mgl32.Vec3{mgl32.Vec3{float32(result.x), float32(result.y), 1}}
		steps = append(i, steps...)
		result = result.p
	}
	return steps
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
	app.n, err = app.Factory().Node()
	if err != nil {
		return nil
	}
	q.SetTexture(t)
	app.n.Attach(q)
	app.n.SetPosition(mgl32.Vec3{0, 0, 1})
	app.renderList = append(app.renderList, app.n)

	// dpad := controller.NewDPad(app.Keyboard())
	// dpad.SetSpeed(3)
	// dpad.SetTarget(n)
	// dpad.SetUpKey(graal.KeyW)
	// dpad.SetDownKey(graal.KeyS)
	// dpad.SetLeftKey(graal.KeyA)
	// dpad.SetRightKey(graal.KeyD)
	// dpad.Track(60.0)

	// app.s = action.NewScheduler()
	// app.disposeList = append(app.disposeList, app.s)
	// go app.s.Start(app.n, 60)
	return nil
}

func (app *Application) moveToPath(v mgl32.Vec3) {
	p := app.n.Position()
	l := app.findPath(uint(p[0]), uint(p[1]), uint(v[0]), uint(v[1]))
	s := action.NewScheduler()
	for _, i := range l {
		dx := float64(i[0] - p[0])
		dy := float64(i[1] - p[1])
		p = i
		l := time.Duration(math.Sqrt(dx*dx+dy*dy)) * time.Millisecond * 500
		s.Add(&action.MoveTo{
			Position: i,
			Duration: l,
		})
	}
	s.Start(app.n, 60)
}

func (app *Application) moveTo(v mgl32.Vec3) {
	p := app.n.Position()
	x, y := float64(p[0]), float64(p[1])
	nx, ny := float64(v[0]), float64(v[1])
	dx, dy := nx-x, ny-y
	x, y = nx, ny
	l := time.Duration(math.Sqrt(dx*dx+dy*dy)*1000) * time.Millisecond
	s := action.NewScheduler()
	s.Add(&action.MoveTo{
		Position: mgl32.Vec3{float32(nx), float32(ny), 1},
		Duration: l,
	})
	go func() {
		s.Start(app.n, 60)
		s.Dispose()
	}()
}

func (app *Application) initMap() error {
	m, err := app.loadMap("assets/tileset.png", 16)
	if err != nil {
		return err
	}
	app.m = m
	app.renderList = append(app.renderList, m)
	return nil
}

func (app *Application) loadMap(m string, s uint) (graal.Tilemap, error) {
	ts, err := app.loadTileset(m)
	if err != nil {
		return nil, err
	}
	defer ts.Release()

	tm, err := app.Factory().Tilemap()
	if err != nil {
		return nil, err
	}
	tm.SetTileset(ts)
	tm.SetSize(s, s)
	for y := uint(0); y < s; y++ {
		for x := uint(0); x < s; x++ {
			tm.SetTile(x, y, 1)
		}
	}
	return tm, nil
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
	eqf := func(a, b float32) bool {
		d := math.Abs(float64(b - a))
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
	).Run(&Application{})
	if err != nil {
		panic(err)
	}
}
