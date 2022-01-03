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

	ps.Model.Position = [3]float64{0, 1.6, 0}
	ps.Model.Velocity = [3]float64{0, 0, 0}
	ps.Model.hardPoints = []mgl64.Vec3{{0, 0, 0}, {-0.0405, -0.49993, 0}, {-2.91, -0.075, 0}}
	ps.SumDT = 0
	ps.StartTime = time.Now()
	ps.LastFrameTime = time.Now()

	ps.Model.Orientation = mgl64.QuatLookAtV(mgl64.Vec3{0, 0, 0}, mgl64.Vec3{0, 0, -1}, mgl64.Vec3{0, 1, 0})

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
	Position mgl64.Vec3
	Velocity mgl64.Vec3
	Mass     float64

	Orientation        mgl64.Quat
	RotationalVelocity mgl64.Vec3
	RotationInertia    mgl64.Mat3

	hardPoints []mgl64.Vec3 //Translations off the origin

}

func (p *PhysicsObject) DoPhysics(dt float64) {
	//dt := frametime.Seconds() //seconds
	//F=m*a
	//mg = ma
	//a = g
	//v += a * dt
	//p += v * dt
	sumfY := 0.0 //p.Mass * g

	for _, v := range p.hardPoints {
		//Collision detection with the ground
		if p.Position.Add(v).Y() < 0 {
			sumfY = 0
			//p.Velocity[1] = 0
		}

	}
	//Linear
	accel := mgl64.Vec3{0, sumfY / p.Mass, 0}
	p.Velocity = p.Velocity.Add(accel.Mul(dt))
	p.Position = p.Position.Add(p.Velocity.Mul(dt))

	//Rotational
	p.Orientation = p.Orientation.Normalize()

	//p.RotationalVelocity = mgl64.Vec3{2 * math.Pi, 0 * math.Pi, 0 * math.Pi}

	// spin = 0.5 *w* *q*
	w := mgl64.Quat{
		W: 0,
		V: [3]float64{p.RotationalVelocity[0], p.RotationalVelocity[1], p.RotationalVelocity[2]},
	}
	q := p.Orientation
	spin := q.Mul(w).Scale(.5)
	p.Orientation = p.Orientation.Add(spin.Scale(dt)).Normalize()

}

//v := p.Orientation.Rotate(mgl64.Vec3{0, 0, 1})
//ang := -math.Atan2(v.Y(), v.Z())
//p.Orientation = mgl64.QuatRotate(ang-0.01, mgl64.Vec3{1, 0, 0}.Normalize())

type CollisionObject struct {
}
