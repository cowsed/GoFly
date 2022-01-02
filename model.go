package main

import (
	graphics "github.com/cowsed/GoFly/Graphics"
	physics "github.com/cowsed/GoFly/Physics"
	"github.com/go-gl/mathgl/mgl32"
)

type Model struct {
	physObj *physics.PhysicsObject
	model3d *graphics.Model
}

func LoadModel(fname string) *Model {
	gfxMod := graphics.MakeModel(fname)

	physObj := &physics.PhysicsObject{
		Position:    [3]float64{0, 10, 0}, //m
		Velocity:    [3]float64{0, 0, 0},  //m/s
		Mass:        10,                   //kg
		Orientation: [3]float64{},
	}
	return &Model{
		physObj: physObj,
		model3d: gfxMod,
	}
}

func (m *Model) DrawModel(projection, view mgl32.Mat4) {
	position := V64toV32(m.physObj.Position)
	orientation := V64toV32(m.physObj.Orientation)
	//position = [3]float32{0, 0, 0}
	m.model3d.DrawModel(projection, view, position, orientation)

}

func V32toV64(v [3]float32) [3]float64 {
	return [3]float64{
		float64(v[0]),
		float64(v[1]),
		float64(v[2]),
	}
}

func V64toV32(v [3]float64) [3]float32 {
	return [3]float32{
		float32(v[0]),
		float32(v[1]),
		float32(v[2]),
	}
}
