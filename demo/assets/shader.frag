#version 330
out vec4 Color;

in vec3 fPosition;
in vec3 fNormal;
in vec2 fTexCoords;
in vec4 fColor;

uniform sampler2D texture1;

void main() {
    Color = texture(texture1, fTexCoords);
}
