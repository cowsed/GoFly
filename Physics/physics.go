package physics

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl64"
)

const g float64 = -9.81 //m/s^2

var AppliedTorque mgl64.Vec3
var AppliedForce mgl64.Vec3

type PhysicsSim struct {
	PhysicsStep   int
	SumDT         float64
	LastFrameTime time.Time
	Model         *PhysicsObject
	Setting       CollisionObject
}

func InitPhysicsContext(mod *PhysicsObject) *PhysicsSim {
	p := PhysicsSim{}
	p.Model = mod
	p.PhysicsStep = 100
	return &p
}
func (ps *PhysicsSim) ResetPhysics() {
	log.Println("Resetting physics")

	ps.Model.Position = [3]float64{0, .5, 0}
	ps.Model.Momentum = [3]float64{0, 0, 0}
	ps.Model.Inertia = ps.Model.Mass * 1.0 / 6.0
	ps.Model.hardPoints = []mgl64.Vec3{
		//Bottom 4
		{-.5, -.5, -.5}, {-.5, -.5, .5}, {.5, -.5, -.5}, {.5, -.5, .5},
		//Top 4
		{-.5, .5, -.5}, {-.5, .5, .5}, {.5, .5, -.5}, {.5, .5, .5},
	}

	//Plane Hardpoints
	//ps.Model.hardPoints = []mgl64.Vec3{{0, 0, 0}, {-0.0405, -0.49993, 0}, {-2.91, -0.075, 0}}
	ps.SumDT = 0
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

	for i := 0; i < ps.PhysicsStep; i++ {
		ps.SumDT += dt / float64(ps.PhysicsStep)

		ps.Model.DoPhysics(ps.SumDT, dt/float64(ps.PhysicsStep))
	}

}

type PhysicsObject struct {
	//Primary
	Position mgl64.Vec3
	Momentum mgl64.Vec3

	//Primary
	Orientation     mgl64.Quat
	AngularMomentum mgl64.Vec3

	//Constant
	Mass    float64
	Inertia float64

	hardPoints []mgl64.Vec3 //Translations off the origin
	Contacting bool
}

func (p *PhysicsObject) ModelSpaceToWorldSpace(offset mgl64.Vec3, centerPosition mgl64.Vec3) mgl64.Vec3 {
	point := offset                                                                                     //Offset in body space
	rotated := p.Orientation.Mat4().Mul4x1(mgl64.Vec4{point[0], point[1], point[2], 1})                 //Offset in rotated body space
	moved := mgl64.Translate3D(centerPosition[0], centerPosition[1], centerPosition[2]).Mul4x1(rotated) //Position of point in world space

	return moved.Vec3()
}
func (p *PhysicsObject) FindCollisions(centerPosition mgl64.Vec3) ([]int, []mgl64.Vec3) { //List of indices to hardpoints array
	hits := []int{}
	pos := []mgl64.Vec3{}
	for i, offset := range p.hardPoints {
		moved := p.ModelSpaceToWorldSpace(offset, centerPosition)
		if moved[1] <= 0 {
			hits = append(hits, i)
			pos = append(pos, moved)
		}
	}

	return hits, pos
}

/*
func (p *PhysicsObject) HandleCollisions(dt float64, nextMomentum, nextPosition mgl64.Vec3) mgl64.Vec3 {
	contactForce := mgl64.Vec3{}
	contactOffset := mgl64.Vec3{}
	p.Contacting = false
	for _, offset := range p.hardPoints {
		point := offset                                                                         //Offset in body space
		rotated := p.Orientation.Mat4().Mul4x1(mgl64.Vec4{point[0], point[1], point[2], 1})     //Offset in rotated body space
		moved := mgl64.Translate3D(p.Position[0], p.Position[1], p.Position[2]).Mul4x1(rotated) //Position of point in world space
		if moved[1] <= 0 {
			p.Contacting = true
			impulseY := 0.0
			//impulse
			//J = integral of F dt
			//J = Favg * t
			//Favg = J/t
			//J = delta P = change in momentum
			//-v to +v for collision with ground
			//J = delta P = mv2-mv1 = m(v2-v1) = m dv
			velocity := p.Momentum.Mul(1 / p.Mass)
			KEBefore := .5 * p.Mass * velocity[1] * velocity[1]
			PEBefore := math.Abs(p.Mass * g * p.Position[1])
			fmt.Printf("Initial:\nKE: %v\nPE: %v\nME: %v\n", KEBefore, PEBefore, KEBefore+PEBefore)
			vy := velocity[1]
			dv := -2. * vy //oppisite direction of velocity * 2 to fully reverse
			impulseY = dv * p.Mass
			forceY := impulseY / dt
			contactForce[1] = forceY
			contactOffset[1] = -(moved[1])
			break //Only handle one force at a time for now

		}

	}
	return contactForce
}
*/
func (p *PhysicsObject) IntegrateLinear(dt float64, force mgl64.Vec3) (mgl64.Vec3, mgl64.Vec3, mgl64.Vec3) {
	//Translational Motion Derivation
	//F=Ma
	//a=F/M
	//Vf = Vi + a*dt
	//Xf = Xi + v*dt
	//P=MV
	//P=Pi+MA*dt
	//P=Pi+M(F/M)*dt // Ms cancel out
	//P=Pi+F*dt

	Momentum := p.Momentum.Add(force.Mul(dt))
	Velocity := Momentum.Mul(1 / p.Mass)
	Position := p.Position.Add(Velocity.Mul(dt))

	return Momentum, Velocity, Position
}
func (p *PhysicsObject) DoPhysics(t, dt float64) {
	gravitationalForce := mgl64.Vec3{0, g * p.Mass, 0}

	force := gravitationalForce                                   //Force in the absense of collisions
	NextMomentum, _, NextPosition := p.IntegrateLinear(dt, force) //Find momentum, position if no collisions

	hits, positions := p.FindCollisions(NextPosition)
	if len(hits) != 0 {
		fmt.Println(hits, positions)
	}
	partialMass := p.Mass / float64(len(hits))

	reactionTorques := make([]mgl64.Vec3, len(hits))
	var reactionForce mgl64.Vec3
	if len(hits) > 0 {
		for i := range hits {
			Vinitial := p.Momentum.Mul(1 / p.Mass)[1]
			MEBefore := .5*p.Mass*Vinitial*Vinitial + p.Mass*p.Position[1]*math.Abs(g)

			p.Position[1] = .5
			PEAfter := (p.Position[1] * math.Abs(g) * p.Mass)
			KEAfter := MEBefore - PEAfter
			if KEAfter < 0 {
				KEAfter = 0
			}
			Vafter := math.Sqrt(2*KEAfter/p.Mass) * .95

			//Impulse = F * t = Ns = Kg * m / s   =  m * dv
			//F*dt = m * dv
			//F= m*dv/dt
			dv := Vafter - Vinitial
			m := partialMass
			F := m * dv / dt
			reactionForce[1] += F
			//p.Momentum = NextMomentum
			//p.Momentum[1] = Vafter * p.Mass
			reactionTorques[i][1] = F
		}
		p.Momentum, _, p.Position = p.IntegrateLinear(dt, force.Add(reactionForce))

	} else {
		p.Momentum = NextMomentum
		p.Position = NextPosition
	}
	//rfmt.Println(reactionTorques)
	//ontactForce := p.HandleCollisions(dt, NextMomentum, NextPosition) //Check if there would be a collision

	//force = force.Add(contactForce) //Add the contact forces in

	//Integrate for real with new forces
	//p.Momentum, _, p.Position = p.IntegrateLinear(dt, force)

	//if p.Contacting {
	//	KEAfter := .5 * p.Mass * Velocity[1] * Velocity[1]
	//	PEAfter := math.Abs(p.Mass * p.Position[1] * g)
	//	fmt.Printf("Final:\nKE: %v\nPE: %v\nME: %v\n", KEAfter, PEAfter, KEAfter+PEAfter)
	//
	//}
	//

	//Rotational
	//torque := AppliedTorque //mgl64.Vec3{.1, 0, 0}
	//T=I alpha
	//alpha = T/I
	//w = wi + alpha * dt
	//I * w = I * wi + alpha*dt * I
	//P = w*I
	//P = Pi + I * alpha * dt
	//alpha = T/I
	//change in P is I * T/I * dt   //Is cancel out
	//P = Pi + T * dt

	//p.AngularMomentum = p.AngularMomentum.Add(torque.Mul(dt))
	//angularVel := p.AngularMomentum.Mul(1 / p.Inertia)
	//p.Orientation = p.Orientation.Normalize()
	//
	////// spin = 0.5 *w* *q*
	//q := mgl64.Quat{
	//	W: 0,
	//	V: [3]float64{angularVel[0], angularVel[1], angularVel[2]},
	//}
	//spin := q.Mul(p.Orientation).Scale(.5)

	//Rotational
	//p.Orientation = p.Orientation.Normalize()

	//t=I alpha
	//Torque := mgl64.Vec3{0, 0, 0}
	//RotationalVelocity := Torque.Mul(1 / p.RotationInertia)
	//

	//spin := q.Mul(w).Scale(.5)
	//p.Orientation = p.Orientation.Add(spin.Scale(dt)).Normalize()

}

//v := p.Orientation.Rotate(mgl64.Vec3{0, 0, 1})
//ang := -math.Atan2(v.Y(), v.Z())
//p.Orientation = mgl64.QuatRotate(ang-0.01, mgl64.Vec3{1, 0, 0}.Normalize())

type CollisionObject struct {
}
