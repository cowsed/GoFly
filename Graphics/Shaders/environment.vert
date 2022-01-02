#version 330

in vec3 vert;

//uniform mat4 viewMatrix;
//uniform mat4 projMatrix;
uniform mat4 MVP;
uniform float SkyBoxScale = 1;


out vec3 fragVert;

void main() {
    fragVert = vert;
	gl_Position = MVP * vec4(vert*SkyBoxScale, 1);
}