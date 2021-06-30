package font

func pow2(value uint32) uint32 {
	value--
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	return value + 1
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}
