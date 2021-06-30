#version 330
out vec4 oColor;

in vec3 fPosition;
in vec3 fNormal;
in vec2 fTexCoords;
in vec4 fColor;

uniform sampler2D texture1;
uniform vec4 color;

void main() {
    oColor = texture(texture1, fTexCoords) * color;
}
