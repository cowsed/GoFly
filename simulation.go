package main

import (
	graphics "github.com/cowsed/GoFly/Graphics"
	physics "github.com/cowsed/GoFly/Physics"
)

type Sim struct {
	mod         *Model
	gfxContext  *graphics.GraphicsContext
	physContext *physics.PhysicsSim
}

func NewSim() *Sim {
	s := Sim{}
	//s.mod = LoadModel("Assets/Models/Simple/final NP v17.obj")
	s.mod = LoadModel(Settings.ModelPath)

	s.gfxContext = graphics.InitGraphicsContext(Settings.EnvironmentPath, Settings.CameraFOV)
	s.physContext = physics.InitPhysicsContext(s.mod.physObj)

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
