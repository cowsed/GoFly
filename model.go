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

func LoadModel(fname string, physObj *physics.PhysicsObject) *Model {
	gfxMod := graphics.MakeModel(fname)

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

func M64toM32(m mgl64.Mat4) mgl32.Mat4 {
	x := [16]float32{
		float32(m[0]),
		float32(m[1]),
		float32(m[2]),
		float32(m[3]),
		float32(m[4]),
		float32(m[5]),
		float32(m[6]),
		float32(m[7]),
		float32(m[8]),
		float32(m[9]),
		float32(m[10]),
		float32(m[11]),
		float32(m[12]),
		float32(m[13]),
		float32(m[14]),
		float32(m[15]),
	}
	return x
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
