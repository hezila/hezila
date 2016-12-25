import (
	"testing"
)

func TestCopy(t *testing.T) {
	va := NewVector(3)
	va.SetValues([]float64{1, 2, 3})
	utils.Expect(t, "1", va.Get(0))
}
