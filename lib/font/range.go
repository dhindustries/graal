package font

type Ranges [][2]rune

func (ranges Ranges) Runes() []rune {
	len := ranges.Length()
	runes := make([]rune, 0, len)
	processed := make(map[rune]bool, len)
	for _, r := range ranges {
		for ch := r[0]; ch < r[1]; ch++ {
			if _, ok := processed[ch]; !ok {
				runes = append(runes, ch)
				processed[ch] = true
			}
		}
	}
	return runes
}

func (ranges Ranges) Length() uint {
	len := uint(0)
	for _, r := range ranges {
		len += uint(r[1]) - uint(r[0])
	}
	return len
}
