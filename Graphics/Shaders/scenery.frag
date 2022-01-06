#version 330

in vec3 fragVert;
out vec4 frag_colour;
uniform float GroundPlaneScale = 1;

uniform vec3 objectPosition = vec3(0,0,0);

float make_grid(vec2 pos){
    float x1 = mod(pos.x,1.0);
    float xMask=step(x1,.5);
    
    float y1 = mod(pos.y,1.0);
    float yMask=step(y1,.5);
 
    return (xMask==0)^^(yMask==0)?1.0:.1;
}



void main() {
    vec3 v = fragVert;
    vec2 uv = v.xz*GroundPlaneScale/2;
    float mask = make_grid(uv);
    vec3 col;
    col = vec3(1,1,1)*mask;

    float shadow = clamp(distance(objectPosition,GroundPlaneScale*fragVert)-.4,0,1);
    shadow*=.8;
    shadow+=.2;
    col*=shadow;
    frag_colour = vec4(col, 1);
}