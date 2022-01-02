package graphics

import (
	_ "embed"
	"fmt"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var ground_points = []float32{
	-1, 0, 1, // Front-top-left
	1, 0, 1, // Front-top-right
	-1, 0, -1, // Front-bottom-left
	1, 0, -1, // Front-bottom-right
}

//go:embed Shaders/scenery.frag
var SceneryFragmentSource string

//go:embed Shaders/scenery.vert
var SceneryVertexSource string

type Scenery struct {
	vao, vbo uint32
	program  uint32
	Scale    float32
}

func (s *Scenery) DrawScenery(projection, view mgl32.Mat4, objPos mgl32.Vec3) {
	gl.UseProgram(s.program)

	mvpMatrixName := "MVP"

	MVP := projection.Mul4(view)
	MVPUniform := gl.GetUniformLocation(s.program, gl.Str(mvpMatrixName+"\x00"))
	gl.UniformMatrix4fv(MVPUniform, 1, false, &MVP[0])

	ScaleUniform := gl.GetUniformLocation(s.program, gl.Str("GroundPlaneScale"+"\x00"))
	gl.Uniform1f(ScaleUniform, s.Scale)

	posUniform := gl.GetUniformLocation(s.program, gl.Str("objectPosition"+"\x00"))
	gl.Uniform3f(posUniform, objPos[0], objPos[1], objPos[2])

	//gl.Disable(gl.CULL_FACE)

	//use the environment things
	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(ground_points)/3))
	//gl.Enable(gl.CULL_FACE)

}

func MakeScenery() *Scenery {
	s := Scenery{}
	s.Scale = 10
	//Make vbo
	gl.GenBuffers(1, &s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(ground_points), gl.Ptr(ground_points), gl.STATIC_DRAW)
	//Make vao
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)

	var err error
	//Build program
	s.program, err = BuildProgram(SceneryFragmentSource, SceneryVertexSource)
	if err != nil {
		panic(fmt.Errorf("error making scenery shader: %v", err))
	}
	return &s
}
