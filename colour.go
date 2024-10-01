package dcolors

import (
	"fmt"
	"image/color"
	"math"
)

// Color represents a color in RGB coordinates. Alpha channel is ignored and always assumed to be equal to 1
//
// This struct maintains internal representation of the color in L*a*b* space which is used to
// calculate perceptual distance between colors. For difference between RGB distance and LAB distance checlk
// https://github.com/lucasb-eyer/go-colorful#comparing-colors
//
// Color implements standard Go color.Color interface
type Color struct {
	rgb [3]uint32
	lab [3]float64
}

// NewColorFromRgb255 creates an instance of Color struct from 8-bit R,G and B channels.
func NewColorFromRgb255(r, g, b uint32) Color {
	return NewColorFromRGBA(r<<8, g<<8, b<<8, 0xFFFF)
}

// NewColorFromRgb creates an instance of Color struct from 16-bit R,G and B channels.
func NewColorFromRgb(r, g, b uint32) Color {
	return NewColorFromRGBA(r, g, b, 0xFFFF)
}

// NewColorFromColor creates and instance of Color struct from the Go's color.Color interface
func NewColorFromColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	if a != 0xffff { // if alpha is not 1, get original rgb values
		r *= 0xffff
		r /= a
		g *= 0xffff
		g /= a
		b *= 0xffff
		b /= a
	}
	return NewColorFromRGBA(r, g, b, a)
}

// NewColorFromRGBA creates a Color from R,G,B and A values. All color channels must be 16-bit.
//
// In contrast with Go color.Color the channel values must NOT be pre-multiplied by alpha.
// Alpha channel is ignored.
func NewColorFromRGBA(r, g, b, _ uint32) Color {
	rf := delinearizeFast(float64(r) / 65535.0)
	gf := delinearizeFast(float64(g) / 65535.0)
	bf := delinearizeFast(float64(b) / 65535.0)
	x, y, z := linearRgbToXyz(rf, gf, bf)
	l, a_, b_ := xyzToLab(x, y, z)

	return Color{rgb: [3]uint32{r, g, b}, lab: [3]float64{l, a_, b_}}
}

// DistanceRgb calculates distance between colors in RGB color space. This is not very useful
func (v Color) DistanceRgb(other Color) float64 {
	dx := v.rgb[0] - other.rgb[0]
	dy := v.rgb[1] - other.rgb[1]
	dz := v.rgb[2] - other.rgb[2]

	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

// Distance calculates the distance between colors using L*a*b* color space
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
