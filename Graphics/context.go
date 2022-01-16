package graphics

import (
	_ "embed"
	"fmt"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

//go:embed Shaders/fullscreenquad.vert
var PPVertSource string

//go:embed Shaders/post_process.frag
var PPFragSource string

//go:embed Shaders/shadow.vert
var ShadowVertSource string

//go:embed Shaders/shadow.frag
var ShadowFragSource string

var points []float32 = []float32{
	-1, 1, 0,
	-1, -1, 0,
	1, -1, 0,

	1, -1, 0,
	1, 1, 0,
	-1, 1, 0,
}

type Camera struct {
	Position mgl32.Vec3
	Lookat   mgl32.Vec3
	FOV      float32
}

type GraphicsContext struct {
	clearColor [4]float32

	Framebuffer                          uint32
	Texture                              uint32
	PPFramebuffer                        uint32 //post processing framebuffer
	PPTexture                            uint32 //post processing texture
	PPProgram                            uint32 //Post processing program
	fullscreenquadVAO, fullscreenquadVBO uint32

	RenderWidth, RenderHeight int32

	LightPosition     mgl32.Vec3
	LightFocus        mgl32.Vec3
	lightSpaceMatrix  mgl32.Mat4
	ShadowFramebuffer uint32
	ShadowTexture     uint32
	ShadowProgram     uint32

	Cam   Camera
	Env   *Environment
	Mod   *Model
	Scene *Model

	Projection, View mgl32.Mat4
}

func InitGraphicsContext(EnvPath, SceneryPath string, camFOV float32) *GraphicsContext {
	gc := GraphicsContext{}
	gc.clearColor = [4]float32{1, 0, 0, 1}
	gc.RenderWidth = 1000
	gc.RenderHeight = 800
	gc.Cam = Camera{
		Position: [3]float32{0, .1, 1},
		Lookat:   [3]float32{0, 0, 0},
		FOV:      camFOV,
	}
	gc.UpdateRenderTargets()

	//Build Post processing program
	var err error
	gc.PPProgram, err = BuildProgram(PPFragSource, PPVertSource)
	if err != nil {
		panic(fmt.Errorf("error building post processing: %v", err))
	}

	gc.CreateShadowLight()

	//Build full screen quad
	gl.GenBuffers(1, &gc.fullscreenquadVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, gc.fullscreenquadVBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &gc.fullscreenquadVAO)
	gl.BindVertexArray(gc.fullscreenquadVAO)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, gc.fullscreenquadVBO)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gc.Env = LoadEnv(EnvPath)

	return &gc
}

const ShadowWidth, SHadowHeight = 1024 * 3, 1024 * 3

func (gfx *GraphicsContext) CreateShadowLight() {
	gfx.LightPosition = mgl32.Vec3{10, 10, 10}
	gfx.LightFocus = mgl32.Vec3{0, 00, 00}

	gl.GenFramebuffers(1, &gfx.ShadowFramebuffer)
	gl.GenTextures(1, &gfx.ShadowTexture)
	gl.BindTexture(gl.TEXTURE_2D, gfx.ShadowTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, ShadowWidth, SHadowHeight, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)

	gl.BindFramebuffer(gl.FRAMEBUFFER, gfx.ShadowFramebuffer)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, gfx.ShadowTexture, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	//Build Post processing program
	var err error
	gfx.ShadowProgram, err = BuildProgram(ShadowFragSource, ShadowVertSource)
	if err != nil {
		panic(fmt.Errorf("error building shadow programming shader: %v", err))
	}

}

func (gfx *GraphicsContext) UpdateRenderTargets() {
	gl.DeleteFramebuffers(1, &gfx.Framebuffer)
	gl.DeleteTextures(1, &gfx.Texture)

	gl.DeleteFramebuffers(1, &gfx.PPFramebuffer)
	gl.DeleteTextures(1, &gfx.PPTexture)

	gl.GenFramebuffers(1, &gfx.Framebuffer)

	//Rendered Texture
	gl.GenTextures(1, &gfx.Texture)
	gl.BindTexture(gl.TEXTURE_2D, gfx.Texture)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gfx.RenderWidth, gfx.RenderHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.BindFramebuffer(gl.FRAMEBUFFER, gfx.Framebuffer)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gfx.Texture, 0)

	//Depth Buffer
	var rbo uint32
	gl.CreateRenderbuffers(1, &rbo)

	gl.BindRenderbuffer(gl.RENDERBUFFER, rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, gfx.RenderWidth, gfx.RenderHeight)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, rbo)
	//Bind back the default framebuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	//Post processing framebuffer
	gl.GenFramebuffers(1, &gfx.PPFramebuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, gfx.PPFramebuffer)

	//Post processing Texture
	gl.GenTextures(1, &gfx.PPTexture)
	gl.BindTexture(gl.TEXTURE_2D, gfx.PPTexture)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gfx.RenderWidth, gfx.RenderHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.BindFramebuffer(gl.FRAMEBUFFER, gfx.PPFramebuffer)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gfx.PPTexture, 0)

	//Bind back the default framebuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

}
func (gfx *GraphicsContext) BeginDraw(objPos mgl32.Vec3) {
	var near, far float32 = 1, 40
	lightProjection := mgl32.Ortho(-10, 10, -10, 10, near, far)
	lightView := mgl32.LookAtV(gfx.LightPosition, gfx.LightFocus, mgl32.Vec3{0, 1, 0})
	gfx.lightSpaceMatrix = lightProjection.Mul4(lightView)

	gl.BindFramebuffer(gl.FRAMEBUFFER, gfx.ShadowFramebuffer)
	gl.Viewport(0, 0, ShadowWidth, SHadowHeight)
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0, 0, 0, 1)

	gfx.RenderScene(true)

	//Calculate perspective and view matrices
	gfx.Projection = mgl32.Perspective(mgl32.DegToRad(gfx.Cam.FOV), float32(gfx.RenderWidth)/float32(gfx.RenderHeight), 0.01, 10000.0)
	gfx.View = mgl32.LookAtV(gfx.Cam.Position, gfx.Cam.Lookat, [3]float32{0, 1, 0})

	//Bind Framebuffer to render
	gl.BindFramebuffer(gl.FRAMEBUFFER, gfx.Framebuffer)
	gl.Viewport(0, 0, gfx.RenderWidth, gfx.RenderHeight)

	gl.ClearColor(gfx.clearColor[0], gfx.clearColor[1], gfx.clearColor[2], gfx.clearColor[3])
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gfx.RenderScene(false)

}

func (gfx *GraphicsContext) DrawModels() {
	gfx.Scene.DrawModel(gfx.Projection, gfx.View, gfx.lightSpaceMatrix, gfx.ShadowTexture)
	gfx.Mod.DrawModel(gfx.Projection, gfx.View, gfx.lightSpaceMatrix, gfx.ShadowTexture)
}

func (gfx *GraphicsContext) RenderScene(IsShadowPass bool) {
	if !IsShadowPass {
		//Draw Environment
		gfx.Env.DrawEnvironment(gfx.Projection, gfx.View)
		return
	}
	//Draw using shadow program
	gl.UseProgram(gfx.ShadowProgram)
	gl.GetUniformLocation(gfx.ShadowProgram, gl.Str("lightSpaceMatrix\x00"))
	gl.UniformMatrix4fv(0, 1, false, &gfx.lightSpaceMatrix[0])

	gfx.Scene.DrawModel(gfx.lightSpaceMatrix, mgl32.Ident4(), gfx.lightSpaceMatrix, 0)

	gfx.Mod.DrawModel(gfx.lightSpaceMatrix, mgl32.Ident4(), gfx.lightSpaceMatrix, 0)
}

func (gfx *GraphicsContext) EndDraw() {
	//Post Process
	//log.Println("PostProcessing")
	gl.BindFramebuffer(gl.FRAMEBUFFER, gfx.PPFramebuffer)
	gl.Viewport(0, 0, gfx.RenderWidth, gfx.RenderHeight)

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindTexture(gl.TEXTURE_2D, gfx.ShadowTexture)
	gl.UseProgram(gfx.PPProgram)
	gl.BindVertexArray(gfx.fullscreenquadVAO)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(points)/3))

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

}
