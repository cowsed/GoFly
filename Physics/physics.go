package physics

import (
	"log"
	"time"

	"github.com/go-gl/mathgl/mgl64"
)

const g float64 = -9.81 //m/s^2

type PhysicsSim struct {
	SumDT         float64
	StartTime     time.Time
	LastFrameTime time.Time
	Model         *PhysicsObject
	Setting       CollisionObject
}

func InitPhysicsContext(mod *PhysicsObject) *PhysicsSim {
	p := PhysicsSim{}
	p.Model = mod
	return &p
}
func (ps *PhysicsSim) ResetPhysics() {
	log.Println("Resetting physics")

	ps.Model.Position = [3]float64{-4, 1, 0}
	ps.Model.Velocity = [3]float64{1, 0, 0}
	ps.SumDT = 0
	ps.StartTime = time.Now()
	ps.LastFrameTime = time.Now()

}

func (ps *PhysicsSim) DoPhysics(paused bool) {
	Delta := time.Since(ps.LastFrameTime)
	ps.LastFrameTime = time.Now()
	if paused {
		return
	}

	dt := Delta.Seconds() //seconds
	ps.SumDT += dt

	ps.Model.DoPhysics(dt)
}

type PhysicsObject struct {
	Position    mgl64.Vec3
	Velocity    mgl64.Vec3
	Mass        float64
	Orientation mgl64.Vec3
	hardPoints  []mgl64.Vec3 //Translations off the origin

}

func (p *PhysicsObject) DoPhysics(dt float64) {
	//dt := frametime.Seconds() //seconds
	//F=m*a
	//mg = ma
	//a = g
	//v += a * dt
	//p += v * dt

	accel := mgl64.Vec3{0, g, 0}
	p.Velocity = p.Velocity.Add(accel.Mul(dt))
	p.Position = p.Position.Add(p.Velocity.Mul(dt))
	if p.Position[1] <= 0 {
		p.Position[1] = 0
		p.Velocity[1] *= -1
	}
}

type CollisionObject struct {
}
