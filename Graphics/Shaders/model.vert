#version 330
in vec3 vert;
in vec3 normal;
in uint material_index;
in uint objID;


uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projMatrix;
uniform mat4 MVP;

uniform mat4 partMatricies[20];


out vec2 fragTexCoord;
out vec3 fragNormal;
out vec3 fragVert;
out vec3 fragWorldPos;

flat out uint fragMatIndex;
void main() {
    //fragTexCoord = vertTexCoord;
    fragMatIndex = material_index;
    fragNormal = normal;//normalize(normal);
    
    mat4 myTrans = partMatricies[objID];
    vec3 ActualVert = (myTrans * vec4(vert,1)).xyz;


    fragVert = ActualVert;
    
    fragWorldPos = (modelMatrix * vec4(ActualVert,1)).xyz;
	
    gl_Position = MVP * vec4(ActualVert, 1);
}