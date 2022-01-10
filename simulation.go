package main

import (
	graphics "github.com/cowsed/GoFly/Graphics"
	plane_physics "github.com/cowsed/GoFly/Physics"
)

type Sim struct {
	mod         *Model
	gfxContext  *graphics.GraphicsContext
	physContext *plane_physics.PhysicsSim
}

func NewSim() *Sim {
	s := Sim{}
	//s.mod = LoadModel("Assets/Models/Simple/final NP v17.obj")

	s.gfxContext = graphics.InitGraphicsContext(Settings.EnvironmentPath, Settings.CameraFOV)
	s.physContext = plane_physics.InitPhysicsContext()

	s.mod = LoadModel(Settings.ModelPath, s.physContext.Model)

	return &s
}
func (s *Sim) DoPhysics(paused bool) {
	s.physContext.DoPhysics(paused)

}

func (s *Sim) Draw() {
	if followModel {
		s.gfxContext.Cam.Lookat = V64toV32(s.mod.physObj.Position)
	}
	s.gfxContext.BeginDraw(V64toV32(s.mod.physObj.Position))

	s.mod.DrawModel(s.gfxContext.Projection, s.gfxContext.View)

	s.gfxContext.EndDraw()

}
