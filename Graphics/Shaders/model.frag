#version 330
in vec2 fragTexCoord;
in vec3 fragNormal;
in vec3 fragVert;
in vec3 fragWorldPos;
flat in uint fragMatIndex;

uniform vec3 lightpos = vec3(-.4,2.9,-5.77);

//uniform vec3 viewpos;
//uniform int ShadeNormal;
//uniform float ambient;
//uniform vec3 MaterialColor;
vec3 lightColor = vec3(1,1,1);
float ambient = .7;
vec3 viewpos = vec3(5,0,1);
uniform vec3 MaterialColors[6];

layout(location = 0) out vec4 frag_colour;

void main() {
    vec3 MaterialColor = MaterialColors[fragMatIndex];
    //vec3 MaterialColor = vec3(1,0,0);
    vec3 lightDir   = normalize(lightpos - fragWorldPos);
    vec3 viewDir    = normalize(viewpos - fragWorldPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);


    vec3 lpos=normalize(lightpos);
    vec3 col = vec3(1,0,0);
    float amt = dot(lpos, fragNormal);
    amt=ambient+(1-ambient)*amt;

    vec3 finalCol=vec3(1,1,1);
   
   finalCol=MaterialColor*amt;
    frag_colour = vec4(finalCol, 1);
}

