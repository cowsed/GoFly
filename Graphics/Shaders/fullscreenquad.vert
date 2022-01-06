#version 330
in vec3 vp;

out vec2 UV;

void main() {
    UV=vp.xy;
    gl_Position = vec4(vp, 1.0);
}






