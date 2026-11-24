package utils

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ClampInside(value, min, max int) (int, bool) {
	if value < min {
		return min, false
	}
	if value > max {
		return max, false
	}
	return value, true
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
