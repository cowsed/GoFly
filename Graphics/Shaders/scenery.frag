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


float line_grid(vec2 uv){
    float gridColorx,gridColory,gridColor;
    float width =.05;
    uv = fract(uv*2);
    gridColorx = abs(uv.x-.5);    
    gridColorx = 1-smoothstep(.5-width,.5+width,gridColorx);
    gridColory = abs(uv.y-.5);    
    gridColory = 1-smoothstep(.5-width,.5+width,gridColory); 
    
    gridColor = gridColorx+gridColory;

    gridColor=smoothstep(1.8,1.85,gridColor);
    return gridColor;
}


void main() {
    vec3 v = fragVert;
    vec2 uv = v.xz*GroundPlaneScale/2;
    float mask = line_grid(uv);
    vec3 col;
    col = vec3(1,1,1)*mask;

    float shadow = clamp(distance(objectPosition,GroundPlaneScale*fragVert)-.4,0,1);
    shadow*=.8;
    shadow+=.2;
    col*=shadow;
    frag_colour = vec4(col, 1);
}