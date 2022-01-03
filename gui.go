package main

import (
	"fmt"
	"time"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"github.com/go-gl/mathgl/mgl64"
)

var lastFrameTime = time.Now()

func MakeUI() {
	dt := float64(time.Since(lastFrameTime).Microseconds()) / 1000
	lastFrameTime = time.Now()

	controls := g.Layout{
		g.Labelf("FPS: %.2f", (1000 / dt)),
		g.Checkbox("Paused", &Paused),
		g.Checkbox("Follow", &followModel),

		g.InputFloat(&Simulation.gfxContext.Scene.Scale).Label("Scale"),
		g.Labelf("sumDT: %v", Simulation.physContext.SumDT),
		g.Custom(func() {

			DragFloat3("Camera Position", (*[3]float32)(&Simulation.gfxContext.Cam.Position), 0.01, -1000, 1000, "%f")
			DragFloat3("Lookat Position", (*[3]float32)(&Simulation.gfxContext.Cam.Lookat), 0.01, -1000, 1000, "%f")
			DragFloat364("Object Position", (*[3]float64)(&Simulation.physContext.Model.Position), 0.01, -1000, 1000, "%f")
			DragFloat364("Object Velocity", (*[3]float64)(&Simulation.physContext.Model.Velocity), 0.01, -1000, 1000, "%f")
			DragFloat364("Rotational Velocity", (*[3]float64)(&Simulation.physContext.Model.RotationalVelocity), 0.01, -1000, 1000, "%f")
			DragQuat64("Orientation", &Simulation.mod.physObj.Orientation, 0.01, -1, 1, "%f")
			imgui.Text(fmt.Sprint(Simulation.physContext.Model.Orientation))

			if imgui.Button("reset physics") {
				Simulation.physContext.ResetPhysics()
			}
			imgui.DragFloatV("FOV", &Simulation.gfxContext.Cam.FOV, 0.1, .5, 179, "%f", 0)
		}),
	}
	render := g.Custom(func() {
		size := imgui.ContentRegionAvail()

		if int32(size.X) != Simulation.gfxContext.RenderWidth || int32(size.Y) != Simulation.gfxContext.RenderHeight {
			Simulation.gfxContext.RenderWidth = int32(size.X)
			Simulation.gfxContext.RenderHeight = int32(size.Y)

			Simulation.gfxContext.UpdateRenderTargets()
			Simulation.Draw()
		} else {
			w := size.X
			var aspect float32 = float32(Simulation.gfxContext.RenderWidth) / float32(Simulation.gfxContext.RenderHeight)
			size = imgui.Vec2{X: w, Y: w / aspect}

		}
		imgui.ImageV(imgui.TextureID(Simulation.gfxContext.Texture),
			size,
			imgui.Vec2{X: 0, Y: 1},
			imgui.Vec2{X: 1, Y: 0},
			imgui.Vec4{X: 1, Y: 1, Z: 1, W: 1},
			imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0},
		)
	})
	if !fullWindow {
		g.SingleWindow().Layout(
			g.SplitLayout(g.DirectionHorizontal, 400, controls, render),
		)
		return
	}
	//If render is maximized
	g.PushWindowPadding(0, 0)
	g.SingleWindow().Layout(render)
	imgui.PopStyleVar()

}

func MakeLoadingUI() {
	g.SingleWindow().Layout(
		g.Label("Please Wait... Loading"),
	)
}

func DragFloat3(label string, vec *[3]float32, speed float32, min, max float32, format string) bool {
	value_changed := false
	size := imgui.CalcItemWidth() / float32(len(vec)+1)
	for i := range vec {
		imgui.PushItemWidth(size)
		id := fmt.Sprintf("%s-%d\n", label, i)
		imgui.PushID(id)
		if i > 0 {
			imgui.SameLine()
		}
		changed := imgui.DragFloatV("", &vec[i], speed, min, max, format, 0)
		value_changed = value_changed || changed
		imgui.PopID()
		imgui.PopItemWidth()
	}

	imgui.SameLine()
	imgui.Text(label)

	return value_changed
}

func DragFloat364(label string, vec64 *[3]float64, speed float32, min, max float32, format string) bool {
	vec := V64toV32(*vec64)

	value_changed := false
	size := imgui.CalcItemWidth() / float32(len(vec)+1)
	for i := range vec {
		imgui.PushItemWidth(size)
		id := fmt.Sprintf("%s-%d\n", label, i)
		imgui.PushID(id)
		if i > 0 {
			imgui.SameLine()
		}
		changed := imgui.DragFloatV("", &vec[i], speed, min, max, format, 0)
		value_changed = value_changed || changed
		imgui.PopID()
		imgui.PopItemWidth()
	}

	imgui.SameLine()
	imgui.Text(label)
	*vec64 = V32toV64(vec)
	return value_changed
}

func DragQuat64(label string, vec64 *mgl64.Quat, speed float32, min, max float32, format string) bool {
	q := Quat64toQuat32(*vec64)
	vec := [4]float32{q.W, q.V.X(), q.V.Y(), q.V.Z()}

	value_changed := false
	size := imgui.CalcItemWidth() / float32(len(vec)+1)
	for i := range vec {
		imgui.PushItemWidth(size)
		id := fmt.Sprintf("%s-%d\n", label, i)
		imgui.PushID(id)
		if i > 0 {
			imgui.SameLine()
		}
		changed := imgui.DragFloatV("", &vec[i], speed, min, max, format, 0)
		value_changed = value_changed || changed
		imgui.PopID()
		imgui.PopItemWidth()
	}

	imgui.SameLine()
	imgui.Text(label)
	*vec64 = mgl64.Quat{
		W: float64(vec[0]),
		V: [3]float64{float64(vec[1]), float64(vec[2]), float64(vec[3])},
	}
	return value_changed
}
