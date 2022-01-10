package plane_physics

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
)

const g float64 = -9.81 //m/s^2

type PhysicsSim struct {
	SumDT         float64
	PhysicsFrames int
	Model         *PhysicsObject
	Gravity       mgl64.Vec3
	//Setting     *bulletphysics.PhysicsObject
}

func InitPhysicsContext() *PhysicsSim {
	p := PhysicsSim{}
	p.PhysicsFrames = 1000
	p.Gravity = mgl64.Vec3{0, g, 0}

	m := PhysicsObject{
		Position:    [3]float64{0, 1, 0},
		Momentum:    [3]float64{0, 0, 0},
		Orientation: mgl64.Quat{},
		Mass:        0,
		contactPoints: []mgl64.Vec3{
			//Bottom 4
			{-.5, -.5, -.5}, {-.5, -.5, .5}, {.5, -.5, -.5}, {.5, -.5, .5},
			//Top 4
			{-.5, .5, -.5}, {-.5, .5, .5}, {.5, .5, -.5}, {.5, .5, .5},
		},
	}

	s := 1.0
	it := mgl64.Mat3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}.Mul(m.Mass * (s * s) / 6)

	m.InertiaTensor = it

	p.Model = &m

	p.ResetPhysics()
	return &p
}

func (ps *PhysicsSim) ResetPhysics() {
	ps.Model.Mass = 1
	ps.Model.Momentum = mgl64.Vec3{}
	ps.Model.Position = mgl64.Vec3{0, 1, 0}
	ps.Model.Orientation = mgl64.QuatRotate(0, mgl64.Vec3{0, 1, 0})
	ps.Model.AngularMomentum = mgl64.Vec3{}

	s := 1.0
	it := mgl64.Mat3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}.Mul(ps.Model.Mass * (s * s) / 6)

	ps.Model.InertiaTensor = it
}

func (ps *PhysicsSim) DoPhysics(paused bool) {
	if paused {
		return
	}

	frameTime := (16 * time.Millisecond).Seconds() //verbosity good
	ps.SumDT += frameTime

	dt := frameTime / float64(ps.PhysicsFrames)
	for i := 0; i < ps.PhysicsFrames; i++ {
		forces := ps.Gravity

		ps.Model.IntegrateLinear(dt, forces)
		lastOrientation := ps.Model.Orientation
		lastAngularMomentum := ps.Model.AngularMomentum
		ps.Model.IntegrateRotational(dt, mgl64.Vec3{})

		collisions := ps.DetectCollisions()
		if len(collisions) == 0 {
			continue
		}
		//If multiple collisions are happening(flat plane on flat plane) collisions can just happen at their center
		//reaction := collisions[0]
		//reaction.CollisionNormal = collisions[0].CollisionNormal
		//for _, r := range collisions {
		//	reaction.CollisionBodyPosition = reaction.CollisionBodyPosition.Add(r.CollisionBodyPosition)
		//}
		//
		//reaction.CollisionBodyPosition = reaction.CollisionBodyPosition.Mul(1.0 / float64(4))
		//Do the actual math
		coeff_restitution := .5

		reactionForce := mgl64.Vec3{}
		reactionTorque := mgl64.Vec3{}

		for _, reaction := range collisions {
			e := coeff_restitution
			v := ps.Model.GetVelocityAtPoint(reaction.CollisionBodyPosition) //ps.Model.Momentum.Mul(1 / ps.Model.Mass)
			n := reaction.CollisionNormal
			r := reaction.CollisionBodyPosition //.Sub(ps.Model.Position)
			inv_mass := 1 / ps.Model.Mass
			inv_it := ps.Model.InertiaTensor.Inv()

			numerator := (v.Mul(-(1 + e))).Dot(n)
			denom := (inv_it.Mul3x1(r.Cross(n))).Cross(r).Dot(n) + inv_mass

			impulse := numerator / denom

			dP := n.Mul(impulse)
			F := dP.Mul(1 / dt)
			reactionForce = reactionForce.Add(F)

			//ps.Model.AngularMomentum = ps.Model.AngularMomentum.Add(r.Cross(n.Mul(j)))

			reactionTorque = reactionTorque.Add(r.Cross(n.Mul(impulse)).Mul(1 / dt))
		}
		//ps.Model.IntegrateLinear(dt, reactionForce)
		ps.Model.Momentum = ps.Model.Momentum.Add(reactionForce.Mul(dt))
		ps.Model.Position = ps.Model.Position.Add(reactionForce.Mul(dt).Mul(dt / ps.Model.Mass))
		//ps.Model.IntegrateRotational(dt, reactionTorque)
		//Do the smae thing as linearly, cant integrate because that adds last position to
		ps.Model.AngularMomentum = lastAngularMomentum
		ps.Model.Orientation = lastOrientation
		ps.Model.IntegrateRotational(dt, reactionTorque)
	}

}
func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func transformVector(mat mgl64.Mat3, v mgl64.Vec3) mgl64.Vec3 {
	return mat.Mul3x1(v)
}

type CollisionResponse struct {
	CollisionNormal       mgl64.Vec3
	CollisionBodyPosition mgl64.Vec3
}

func (ps *PhysicsSim) DetectCollisions() []CollisionResponse {
	reactions := []CollisionResponse{}
	for i := range ps.Model.contactPoints {
		world_space := ps.Model.ModelSpaceToWorldSpace(ps.Model.contactPoints[i], ps.Model.Position)
		if world_space[1] < 0 {
			reactions = append(reactions, CollisionResponse{
				CollisionNormal:       [3]float64{0, 1, 0},
				CollisionBodyPosition: ps.Model.contactPoints[i],
			})
		}
	}
	return reactions
}

type PhysicsObject struct {
	Position mgl64.Vec3
	Momentum mgl64.Vec3

	Orientation     mgl64.Quat
	AngularMomentum mgl64.Vec3

	Mass          float64
	InertiaTensor mgl64.Mat3

	contactPoints []mgl64.Vec3
}

func (p *PhysicsObject) GetVelocityAtPoint(point mgl64.Vec3) mgl64.Vec3 {
	inverseInertiaTensor := p.InertiaTensor.Inv()
	angularVelocity := transformVector(inverseInertiaTensor, p.AngularMomentum)

	linearVel := p.Momentum.Mul(1 / p.Mass)
	vel := linearVel.Add(angularVelocity.Cross(point.Sub(p.Position)))
	return vel
}

func (p *PhysicsObject) IntegrateLinear(dt float64, forces mgl64.Vec3) {
	//V_{n+1} = V_n + A*dt
	//A = F/m
	//P = Vm
	//P_{n+1} = P_{n} + m*A*dt
	//P_{n+1} = P_{n} + m*F/m*dt
	//P_{n+1} = P_{n} + F*dt

	p.Momentum = p.Momentum.Add(forces.Mul(dt))
	p.Position = p.Position.Add(p.Momentum.Mul(1 / p.Mass).Mul(dt))
}
func (p *PhysicsObject) IntegrateRotational(dt float64, torques mgl64.Vec3) {
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

	p.AngularMomentum = p.AngularMomentum.Add(torques.Mul(dt))

	inverseInertiaTensor := p.InertiaTensor.Inv()
	angularVel := inverseInertiaTensor.Mul3x1(p.AngularMomentum)
	p.Orientation = p.Orientation.Normalize()

	//// spin = 0.5 *w* *q*
	q := mgl64.Quat{
		W: 0,
		V: [3]float64{angularVel[0], angularVel[1], angularVel[2]},
	}
	spin := q.Mul(p.Orientation).Scale(.5)
	p.Orientation = p.Orientation.Add(spin.Scale(dt)).Normalize()

}

func (p *PhysicsObject) ModelSpaceToWorldSpace(ObjectSpace mgl64.Vec3, centerPosition mgl64.Vec3) mgl64.Vec3 {
	point := ObjectSpace                                                                                //Offset in body space
	rotated := p.Orientation.Mat4().Mul4x1(mgl64.Vec4{point[0], point[1], point[2], 1})                 //Offset in rotated body space
	moved := mgl64.Translate3D(centerPosition[0], centerPosition[1], centerPosition[2]).Mul4x1(rotated) //Position of point in world space
	return moved.Vec3()
}

func ExtractPosition(m mgl64.Mat4) mgl64.Vec3 {

	//scalingFactor := math.Sqrt(m.At(0, 0)*m.At(0, 0) + m.At(0, 1)*m.At(0, 1) + m.At(0, 2)*m.At(0, 2))

	translation := mgl64.Vec3{
		m.At(0, 3),
		m.At(1, 3),
		m.At(2, 3),
	}
	return translation
}
