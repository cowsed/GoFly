package graphics

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

//Builds the opengl program
func BuildProgram(FragSrc, VertSrc string) (uint32, error) {

	Program := gl.CreateProgram()

	var vertexShader, fragmentShader uint32
	var err error

	//Compile Vertex Shader
	vertexShader, err = compileShader(VertSrc+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		err = fmt.Errorf("vertex shader error: %s", err)
		return 0, err
	}
	gl.AttachShader(Program, vertexShader)

	//Compile Fragment Shader
	fragmentShader, err = compileShader(FragSrc+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		err = fmt.Errorf("fragment shader error: %s", err)
		return 0, err
	}
	gl.AttachShader(Program, fragmentShader)

	gl.LinkProgram(Program)

	//Release programs
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	//Check Link Errors
	var isLinked int32
	gl.GetProgramiv(Program, gl.LINK_STATUS, &isLinked)
	if isLinked == gl.FALSE {
		var maxLength int32
		gl.GetProgramiv(fragmentShader, gl.INFO_LOG_LENGTH, &maxLength)

		infoLog := make([]uint8, maxLength+1) //[bufSize]uint8{}
		gl.GetShaderInfoLog(fragmentShader, maxLength, &maxLength, &infoLog[0])

		return 0, fmt.Errorf("error linking: %s", string(infoLog))

	}
	return Program, nil
}

//Compiles shaders
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile:\nLog:\n%v", log[:len(log)-2])
	}
	return shader, nil
}
