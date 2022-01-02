package graphics

import (
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	_ "github.com/mdouchement/hdr/codec/rgbe"
)

//go:embed Shaders/environment.frag
var EnvironmentFragmentSource string

//go:embed Shaders/environment.vert
var EnvironmentVertexSource string

var cube_strip = []float32{
	-1, 1, 1, // Front-top-left
	1, 1, 1, // Front-top-right
	-1, -1, 1, // Front-bottom-left
	1, -1, 1, // Front-bottom-right
	1, -1, -1, // Back-bottom-right
	1, 1, 1, // Front-top-right
	1, 1, -1, // Back-top-right
	-1, 1, 1, // Front-top-left
	-1, 1, -1, // Back-top-left
	-1, -1, 1, // Front-bottom-left
	-1, -1, -1, // Back-bottom-left
	1, -1., -1, // Back-bottom-right
	-1, 1., -1, // Back-top-left
	1, 1., -1, // Back-top-right
}

type Environment struct {
	cubemap  uint32
	vao, vbo uint32
	program  uint32
	Scale    float32
}

func LoadEnv(fpath string) *Environment {
	e := Environment{}
	e.Scale = 3000
	e.cubemap = loadCubemaps(fpath)
	//Build Object

	//Make vbo
	gl.GenBuffers(1, &e.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, e.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(cube_strip), gl.Ptr(cube_strip), gl.STATIC_DRAW)
	//Make vao
	gl.GenVertexArrays(1, &e.vao)
	gl.BindVertexArray(e.vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)

	var err error
	//Build program
	e.program, err = BuildProgram(EnvironmentFragmentSource, EnvironmentVertexSource)
	if err != nil {
		panic(fmt.Errorf("error making environment shader: %v", err))
	}

	return &e
}

func (env *Environment) DrawEnvironment(projection, view mgl32.Mat4) {
	gl.UseProgram(env.program)

	mvpMatrixName := "MVP"

	MVP := projection.Mul4(view)
	MVPUniform := gl.GetUniformLocation(env.program, gl.Str(mvpMatrixName+"\x00"))
	gl.UniformMatrix4fv(MVPUniform, 1, false, &MVP[0])

	ScaleUniform := gl.GetUniformLocation(env.program, gl.Str("SkyBoxScale"+"\x00"))
	gl.Uniform1f(ScaleUniform, env.Scale)

	gl.Disable(gl.CULL_FACE)
	//use the environment things
	gl.BindVertexArray(env.vao)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, env.cubemap)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(cube_strip)/3))

	gl.Enable(gl.CULL_FACE)

}

func loadCubemaps(fpath string) uint32 {
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, textureID)
	var sides = []string{
		"right.png",
		"left.png",
		"top.png",
		"bottom.png",
		"front.png",
		"back.png",
	}
	for i := uint32(0); i < 6; i++ {
		f, err := os.Open(fpath + sides[i])
		check(err)
		img, err := png.Decode(f)
		if err == nil {

			var dataImg *image.RGBA = image.NewRGBA(img.Bounds())
			draw.Draw(dataImg, img.Bounds(), img, image.Pt(0, 0), draw.Src)
			var width, height int32 = int32(dataImg.Rect.Dx()), int32(dataImg.Rect.Dy())

			gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(dataImg.Pix))

		} else {
			panic(err)
		}
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	return textureID

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
