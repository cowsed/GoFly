package main

import (
	graphics "github.com/cowsed/GoFly/Graphics"
	physics "github.com/cowsed/GoFly/Physics"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Model struct {
	physObj *physics.PhysicsObject
	model3d *graphics.Model
}

func LoadModel(fname string) *Model {
	gfxMod := graphics.MakeModel(fname)

	physObj := &physics.PhysicsObject{
		Position:    [3]float64{0, 10, 0}, //m
		Momentum:    [3]float64{0, 0, 0},  //m/s
		Mass:        1,                    //kg
		Orientation: mgl64.Quat{},
	}
	return &Model{
		physObj: physObj,
		model3d: gfxMod,
	}
}

func Quat64toQuat32(q mgl64.Quat) mgl32.Quat {
	return mgl32.Quat{
		W: float32(q.W),
		V: V64toV32(q.V),
	}
}

func (m *Model) DrawModel(projection, view mgl32.Mat4) {
	position := V64toV32(m.physObj.Position)
	orientation := Quat64toQuat32(m.physObj.Orientation)
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
