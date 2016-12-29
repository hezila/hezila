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

func maxUInt(x, y uint) uint {
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

func minUInt(x, y uint) uint {
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

type box interface{}

func countBoxes(start, cap uint) chan box {
	ints := make(chan box)
	go func() {
		for i := start; i < cap; i++ {
			ints <- i
		}
		close(ints)
	}()
	return ints
}

func parFor(inputs <-chan box, foo func(i box)) (wait func()) {
	n := runtime.GOMAXPROCS(0)
	block := make(chan bool, n)
	var j uint
	for j = 0; j < n; j++ {
		go func() {
			for {
				i, ok := <-inputs
				if !ok {
					break
				}
				foo(i)
			}
			block <- true
		}()
	}
	wait = func() {
		var i uint
		for i = 0; i < n; i++ {
			<-block
		}
	}
	return
}
