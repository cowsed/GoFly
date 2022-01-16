#version 330
in vec3 fragNormal;
in vec3 fragVert;
in vec3 fragWorldPos;
in vec4 FragPosLightSpace;

flat in uint fragMatIndex;

uniform vec3 lightpos = vec3(100,100,0);
uniform vec3 MaterialColors[6];

uniform mat4 lightSpaceMatrix;
uniform sampler2D ShadowMap;

uniform float shadowBias = 0.0006;

float ambient = .5;
vec3 viewpos = vec3(5,0,1);

layout(location = 0) out vec4 frag_colour;


float ShadowCalc(){
    vec3 pos = FragPosLightSpace.xyz * .5 + .5;
    vec2 lightTexCoords = pos.xy;
    if (lightTexCoords.x<0 || lightTexCoords.x>1 || lightTexCoords.y<0 || lightTexCoords.y>1){
        return 1.0;
    }
    float depth = texture(ShadowMap,lightTexCoords).r;

    return (depth+shadowBias) < pos.z ? 0.0 : 1.0;

}

void main() {
    vec3 MaterialColor = MaterialColors[fragMatIndex];
    //vec3 MaterialColor = vec3(1,0,0);
    vec3 lightDir   = normalize(lightpos - fragWorldPos);
    vec3 viewDir    = normalize(viewpos - fragWorldPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);


    vec3 lpos=normalize(lightpos);
    vec3 col = vec3(1,0,0);
    float amt = dot(lpos, fragNormal);
    
    vec3 finalCol=vec3(1,1,1);
   
    float shadow = ShadowCalc();
    finalCol = (shadow * amt + ambient) * MaterialColor;

    frag_colour = vec4(finalCol, 1);
}

