package text

import (
	"github.com/dhindustries/graal"
)

type Character struct {
	Texture graal.Texture
	Source  [4]float64
	Size    [2]uint64
	Bearing [2]uint64
	Advance uint64
}
