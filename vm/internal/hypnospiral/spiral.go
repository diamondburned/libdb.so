package hypnospiral

import (
	"cmp"
	"image"
	"image/color"
	"math"
)

func clamp[T cmp.Ordered](v, vmin, vmax T) T {
	return min(max(v, vmin), vmax)
}

type Parameters struct {
	SpinSpeed float64
	Zoom      float64
	Blur      float64
}

var DefaultParameters = Parameters{
	SpinSpeed: 3.6,
	Zoom:      1.2,
	Blur:      0.4,
}

// Draw renders the hypnotic spiral pattern onto the provided grayscale image
// buffer.
func Draw(buf *image.Gray, p Parameters, t float64) {
	angle := t * (2 * math.Pi / 5) * p.SpinSpeed

	hw := float64(buf.Rect.Dx()) / 2
	hh := float64(buf.Rect.Dy()) / 2

	for y := range buf.Rect.Dy() {
		for x := range buf.Rect.Dx() {
			dx := float64(x) - hw
			dy := float64(y) - hh

			r := math.Hypot(dx, dy)
			v := math.Sin(math.Atan2(dy, dx) + math.Sqrt(r)*2.5*p.Zoom + angle)

			// Normalize sine wave to [0, 1]
			v = (v + 1) / 2

			// Use the blur parameter to ease the transition smoothly.
			// A higher blur value makes the transition softer and more blurred.
			v = clamp((v-0.5)/p.Blur+0.5, 0.0, 1.0)
			v = v * v * (3.0 - 2.0*v)

			buf.SetGray(x, y, color.Gray{Y: uint8(v * 255)})
		}
	}
}
