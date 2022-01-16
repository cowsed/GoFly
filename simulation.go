package main

import (
	"fmt"

	graphics "github.com/cowsed/GoFly/Graphics"
	plane_physics "github.com/cowsed/GoFly/Physics"
	"github.com/go-gl/mathgl/mgl32"
)

type Sim struct {
	mod   *Model
	scene *Model

	gfxContext  *graphics.GraphicsContext
	physContext *plane_physics.PhysicsSim
}

func NewSim() *Sim {
	s := Sim{}
	//s.mod = LoadModel("Assets/Models/Simple/final NP v17.obj")

	s.gfxContext = graphics.InitGraphicsContext(Settings.EnvironmentPath, Settings.SceneryPath, Settings.CameraFOV)
	s.physContext = plane_physics.InitPhysicsContext()
	s.mod = LoadModel(Settings.ModelPath, s.physContext.Model)
	s.scene = LoadModel(Settings.SceneryPath, nil)
	s.scene.model3d.ModelMatrix = mgl32.Ident4()
	s.gfxContext.Mod = s.mod.model3d
	s.gfxContext.Scene = s.scene.model3d

	fmt.Println("SCENE", s.scene.model3d)
	return &s
}
func (s *Sim) DoPhysics(paused bool) {
	s.physContext.DoPhysics(paused)

}

func (s *Sim) Draw() {
	if followModel {
		s.gfxContext.Cam.Lookat = V64toV32(s.mod.physObj.Position)
	}
	s.mod.ApplyPhysics()

	s.gfxContext.BeginDraw(V64toV32(s.mod.physObj.Position))
	s.gfxContext.DrawModels()

	s.gfxContext.EndDraw()

}
