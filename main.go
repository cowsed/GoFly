package main

import (
	"log"
	"time"

	g "github.com/AllenDang/giu"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var Simulation *Sim
var frame = 0
var Paused = true
var fullWindow = false
var followModel = false

func refresh() {
	ticker := time.NewTicker(time.Millisecond * 16)
	for {
		g.Update()
		<-ticker.C
	}
}

var wnd *g.MasterWindow
var Settings Config

func main() {
	Settings = LoadConfig()
	wnd = g.NewMasterWindow("GoFly", 1800, 980, 0)

	//Update window every 16ms
	go refresh()

	//Run application
	wnd.Run(loop)

}

func loop() {
	if frame == 0 {
		log.Println("Making ui")
		MakeLoadingUI()
		frame = 1
		return
	}

	if frame == 1 {
		log.Println("Making sim")
		Simulation = NewSim()
		//Simulation.physContext.Run
		Simulation.physContext.ResetPhysics()
		Paused = false
	}
	Simulation.Draw()

	MakeUI()

	Simulation.DoPhysics(Paused)

	if !Paused {
		frame++
	}

	if g.IsKeyPressed(g.KeyEscape) {
		Paused = !Paused
	}
	if g.IsKeyPressed(g.KeyR) {
		Simulation.physContext.ResetPhysics()
	}
	if g.IsKeyPressed(g.KeyF) {
		fullWindow = !fullWindow
	}
	if g.IsKeyPressed(g.KeyG) {
		//wnd.SetSize(1920, 1080)
		//	wnd.SetPos(0, 0)
		glfw.GetCurrentContext().Maximize()
	}

}
