// The code in this file is extracted from go-colorful library
// https://github.com/lucasb-eyer/go-colorful
// This file is licensed under MIT license
// SPDX-License-Identifier: MIT

package dominant_colors

import "math"

func sq(v float64) float64 {
	return v * v
}
func delinearize(v float64) float64 {
	if v <= 0.0031308 {
		return 12.92 * v
	}
	return 1.055*math.Pow(v, 1.0/2.4) - 0.055
}

func delinearizeFast(v float64) float64 {
	// This function (fractional root) is much harder to linearize, so we need to split.
	if v > 0.2 {
		v1 := v - 0.6
		v2 := v1 * v1
		v3 := v2 * v1
		v4 := v2 * v2
		v5 := v3 * v2
		return 0.442430344268235 + 0.592178981271708*v - 0.287864782562636*v2 + 0.253214392068985*v3 - 0.272557158129811*v4 + 0.325554383321718*v5
	} else if v > 0.03 {
		v1 := v - 0.115
		v2 := v1 * v1
		v3 := v2 * v1
		v4 := v2 * v2
		v5 := v3 * v2
		return 0.194915592891669 + 1.55227076330229*v - 3.93691860257828*v2 + 18.0679839248761*v3 - 101.468750302746*v4 + 632.341487393927*v5
	} else {
		v1 := v - 0.015
		v2 := v1 * v1
		v3 := v2 * v1
		v4 := v2 * v2
		v5 := v3 * v2
		// You can clearly see from the involved constants that the low-end is highly nonlinear.
		return 0.0519565234928877 + 5.09316778537561*v - 99.0338180489702*v2 + 3484.52322764895*v3 - 150028.083412663*v4 + 7168008.42971613*v5
	}
}

func linearRgbToXyz(r, g, b float64) (x, y, z float64) {
	x = 0.41239079926595948*r + 0.35758433938387796*g + 0.18048078840183429*b
	y = 0.21263900587151036*r + 0.71516867876775593*g + 0.072192315360733715*b
	z = 0.019330818715591851*r + 0.11919477979462599*g + 0.95053215224966058*b
	return
}

// This is the default reference white point.
var wref = [3]float64{0.95047, 1.00000, 1.08883}

// func XyzToLabWhiteRef(x, y, z float64, wref [3]float64) (l, a, b float64) {
func xyzToLab(x, y, z float64) (l, a, b float64) {
	fy := labF(y / wref[1])
	l = 1.16*fy - 0.16
	a = 5.0 * (labF(x/wref[0]) - fy)
	b = 2.0 * (fy - labF(z/wref[2]))
	return
}

func labF(t float64) float64 {
	if t > 6.0/29.0*6.0/29.0*6.0/29.0 {
		return math.Cbrt(t)
	}
	return t/3.0*29.0/6.0*29.0/6.0 + 4.0/29.0
}
