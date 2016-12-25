package math

func fibonacci() func() int {
	var x, y int = -1, 1
	return func() int {
		x, y = y, x+y
		return y
	}
}

func max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func sum(a []float64) (s float64) {
	for _, v := range a {
		s += v
	}
	return
}

func product(a []float64) float64 {
	p := float64(1)
	for _, v := range a {
		p *= v
	}
	return p
}
