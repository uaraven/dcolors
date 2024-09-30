package dominant_colors

import (
	"fmt"
	"image/color"
	"math"
)

type Color struct {
	rgb [3]uint32
	lab [3]float64
}

func NewColorFromRgb255(r, g, b uint32) Color {
	return NewColourFromRGBA(r<<8, g<<8, b<<8, 0xFFFF)
}

func NewColorFromRgb(r, g, b uint32) Color {
	return NewColourFromRGBA(r, g, b, 0xFFFF)
}

func NewColourFromColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	return NewColourFromRGBA(r, g, b, a)
}

func NewColourFromRGBA(r, g, b, a uint32) Color {
	rf := delinearizeFast(float64(r) / 65535.0)
	gf := delinearizeFast(float64(g) / 65535.0)
	bf := delinearizeFast(float64(b) / 65535.0)
	x, y, z := linearRgbToXyz(rf, gf, bf)
	l, a_, b_ := xyzToLab(x, y, z)

	return Color{rgb: [3]uint32{r, g, b}, lab: [3]float64{l, a_, b_}}
}

func (v Color) DistanceRgb(other Color) float64 {
	dx := v.rgb[0] - other.rgb[0]
	dy := v.rgb[1] - other.rgb[1]
	dz := v.rgb[2] - other.rgb[2]

	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

func (v Color) Distance(c2 Color) float64 {
	l1, a1, b1 := v.lab[0], v.lab[1], v.lab[2]
	l2, a2, b2 := c2.lab[0], c2.lab[1], c2.lab[2]
	return math.Sqrt(sq(l1-l2) + sq(a1-a2) + sq(b1-b2))
}

func (v Color) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", v.rgb[0]>>8, v.rgb[1]>>8, v.rgb[2]>>8)
}

func (v Color) String() string {
	return v.Hex()
}

func (v Color) AsColor() color.Color {
	return color.RGBA64{R: uint16(v.rgb[0]), G: uint16(v.rgb[1]), B: uint16(v.rgb[2]), A: 0xFFFF}
}

func (v Color) RGBA() (r, g, b, a uint32) {
	return v.rgb[0], v.rgb[1], v.rgb[2], 0xffff
}
