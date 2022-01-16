#version 330
in vec2 UV;
uniform sampler2D outputImage;

layout (location = 0) out vec4 frag_color;

void main(){
    vec3 col = vec3(1,0,0);
    vec2 texCoords = UV/2 +.5;
    col=vec3(texture(outputImage,texCoords).r);
 
    //Gamma correct
    float gamma = 1;//2.2;
    col = pow(col, vec3(1/gamma));
    frag_color = vec4(col,1);
}