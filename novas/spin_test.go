package novas

import (
	"math"
	"testing"
)

func TestSpin(t *testing.T) {
	angle := 33.0
	pos1 := make([]float64, 3)
	pos2 := make([]float64, 3)
	exppos2 := make([]float64, 3)
	exppos2[0] = 0.515266
	exppos2[1] = 0.3027273701968841
	exppos2[2] = 0.8017837257372732

	norm := math.Sqrt(1.0 + 4.0 + 9.0)
	pos1[0] = 1.0 / norm
	pos1[1] = 2.0 / norm
	pos1[2] = 3.0 / norm
	spin(angle, pos1, pos2)
	for idx := 0; idx < 3; idx++ {
		if err := checkFloat(pos2[idx], exppos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
	spin(angle, pos1, pos2)
	for idx := 0; idx < 3; idx++ {
		if err := checkFloat(pos2[idx], exppos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
}
