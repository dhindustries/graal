package main

import (
	"log"
	"math"
	"os"

	"github.com/dhindustries/graal"
	"github.com/dhindustries/graal/components"
	"github.com/dhindustries/graal/glfw"
	"github.com/dhindustries/graal/opengl"
)

type Actor struct {
	components.Shaped
	components.Transformed
}

type Application struct {
	graal.BaseApplication
	renderList []interface{}
	actor      *Actor
}

func (app *Application) Prepare() {
	shape := app.Factory().Quad(-1, -1, 1, 1)
	texture, err := app.Resources().LoadTexture("assets/flower.jpg")
	if err != nil {
		app.Log(err)
	} else {
		defer texture.Release()
		shape.SetTexture(texture)
	}
	app.actor = &Actor{}
	app.actor.SetShape(shape)
	app.actor.SetPosition(graal.Vec3{0, 0, 0})
	app.actor.SetScale(graal.Vec3{1, 1, 1})

	camera := graal.OrhtoCamera{}
	camera.SetNear(-1)
	camera.SetFar(1)
	camera.SetPosition(graal.Vec3{0, 0, 0})
	camera.SetLookAt(graal.Vec3{0, 0, 1})
	camera.SetUp(graal.Vec3{0, 1, 0})
	camera.SetViewport(0, 600, 800, 0)

	// app.Renderer().Use(&camera)
}

func (app *Application) Update(dt float32) {
	shape := app.actor.Shape()
	if app.Keyboard().IsPressed(graal.KeyEscape) {
		app.Close()
	}
	cs := 0.25 * dt
	ps := 0.1 * dt
	if app.Keyboard().IsDown(graal.KeyQ) {
		c := shape.Color()
		c[0] = float32(math.Min(float64(c[0]+cs), 1))
		shape.SetColor(c)
	}
	if app.Keyboard().IsDown(graal.KeyA) {
		c := shape.Color()
		c[0] = float32(math.Max(float64(c[0]-cs), 0))
		shape.SetColor(c)
	}
	if app.Keyboard().IsDown(graal.KeyW) {
		c := shape.Color()
		c[1] = float32(math.Min(float64(c[1]+cs), 1))
		shape.SetColor(c)
	}
	if app.Keyboard().IsDown(graal.KeyS) {
		c := shape.Color()
		c[1] = float32(math.Max(float64(c[1]-cs), 0))
		shape.SetColor(c)
	}
	if app.Keyboard().IsDown(graal.KeyE) {
		c := shape.Color()
		c[2] = float32(math.Min(float64(c[2]+cs), 1))
		shape.SetColor(c)
	}
	if app.Keyboard().IsDown(graal.KeyD) {
		c := shape.Color()
		c[2] = float32(math.Max(float64(c[2]-cs), 0))
		shape.SetColor(c)
	}

	if app.Keyboard().IsDown(graal.KeyLeft) {
		p := app.actor.Position()
		p[0] = p[0] - ps
		app.actor.SetPosition(p)
	}

	if app.Keyboard().IsDown(graal.KeyRight) {
		p := app.actor.Position()
		p[0] = p[0] + ps
		app.actor.SetPosition(p)
	}
	if app.Keyboard().IsDown(graal.KeyUp) {
		p := app.actor.Position()
		p[1] = p[1] - ps
		app.actor.SetPosition(p)
	}
	if app.Keyboard().IsDown(graal.KeyDown) {
		p := app.actor.Position()
		p[1] = p[1] + ps
		app.actor.SetPosition(p)
	}

	r := app.actor.Rotation()
	r[1] += 0.1 * dt
	r[0] += 0.2 * dt
	app.actor.SetRotation(r)
}

func (app *Application) Render() {
	app.Renderer().Render(app.actor)
	// for _, item := range app.renderList {
	// 	app.Renderer().Render(item)
	// }
}

func (app *Application) Dispose() {
	app.actor.Shape().Dispose()
}

func main() {
	engine := graal.Engine{}
	engine.Window = glfw.NewWindow(800, 600, "App")
	engine.Graphics = &opengl.Graphics{}
	engine.Input = &glfw.Input{}
	engine.Logger = log.New(os.Stderr, "", 0)
	engine.Run(&Application{})
}
