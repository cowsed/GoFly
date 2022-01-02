#version 330
uniform samplerCube envTex;

in vec3 fragVert;
out vec4 frag_colour;

void main() {
    vec3 col = normalize(fragVert);
    vec3 v = fragVert;

    vec3 texcol = texture(envTex,v).xyz;
    
    col=mix(col,texcol,1);
    frag_colour = vec4(col, 1);
}