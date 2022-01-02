#version 330

in vec3 vert;

uniform mat4 MVP;
uniform float GroundPlaneScale = 1;


out vec3 fragVert;

void main() {
    fragVert = vert;
	gl_Position = MVP * vec4(vert*GroundPlaneScale, 1);
}