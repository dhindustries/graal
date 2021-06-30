package freetype

type Range [2]rune
type Ranges []Range

func (r Ranges) Length() uint32 {
	len := uint32(0)
	for _, v := range r {
		len += uint32(v[1]) - uint32(v[0])
	}
	return len
}
