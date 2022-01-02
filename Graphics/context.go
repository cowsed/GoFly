package graphics

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position mgl32.Vec3
	Lookat   mgl32.Vec3
	FOV      float32
}

type GraphicsContext struct {
	Cam                       Camera
	clearColor                [4]float32
	Framebuffer               uint32
	Texture                   uint32
	Scene                     *Scenery
	Env                       *Environment
	Mod                       *Model
	RenderWidth, RenderHeight int32
	Projection, View          mgl32.Mat4
}

func InitGraphicsContext(EnvPath string, camFOV float32) *GraphicsContext {
	gc := GraphicsContext{}
	gc.clearColor = [4]float32{1, 0, 0, 1}
	gc.RenderWidth = 1000
	gc.RenderHeight = 800
	gc.Cam = Camera{
		Position: [3]float32{0, 1.6, 5},
		Lookat:   [3]float32{0, 0, 0},
		FOV:      camFOV,
	}
	gc.UpdateRenderTargets()
	LoadACFile("Assets/Planes/allegro.ac")

	gc.Env = LoadEnv(EnvPath)
	gc.Scene = MakeScenery()

	return &gc
}
func (gfx *GraphicsContext) UpdateRenderTargets() {
	gl.DeleteFramebuffers(1, &gfx.Framebuffer)
	gl.DeleteTextures(1, &gfx.Texture)
	gl.GenFramebuffers(1, &gfx.Framebuffer)

	//Rendered Texture
	gl.GenTextures(1, &gfx.Texture)
	gl.BindTexture(gl.TEXTURE_2D, gfx.Texture)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gfx.RenderWidth, gfx.RenderHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.BindFramebuffer(gl.FRAMEBUFFER, gfx.Framebuffer)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gfx.Texture, 0)

	var rbo uint32
	gl.CreateRenderbuffers(1, &rbo)

	gl.BindRenderbuffer(gl.RENDERBUFFER, rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, gfx.RenderWidth, gfx.RenderHeight)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, rbo)

	//Bind back the default framebuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

}
func (gfx *GraphicsContext) BeginDraw(objPos mgl32.Vec3) {
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

	//Draw Environment
	gfx.Env.DrawEnvironment(gfx.Projection, gfx.View)
	//Draw Scenery
	gfx.Scene.DrawScenery(gfx.Projection, gfx.View, objPos)

}
func (gfx *GraphicsContext) EndDraw() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

}
