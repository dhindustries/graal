package text

import "github.com/dhindustries/graal"

type Font interface {
	graal.Handle
	Char(ch rune) *Character
}
