#version 330
layout(location=0) in vec3 vPosition;
layout(location=1) in vec3 vNormal;
layout(location=2) in vec2 vTexCoords;
layout(location=3) in vec4 vColor;

out vec3 fPosition;
out vec3 fNormal;
out vec2 fTexCoords;
out vec4 fColor;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
    fPosition = vPosition;
    fNormal = vNormal;
    fTexCoords = vTexCoords;
    fColor = vColor;
    gl_Position = projection * view * model * vec4(vPosition, 1);
}
