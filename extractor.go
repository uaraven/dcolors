package dominant_colors

import "image"

type Options struct {
	SamplingInterval uint
	ExactMatch       bool
}

func ExtractDominantColors(image image.Image, number int, options *Options) []Color {
	var sampling uint
	exactMatch := false
	if options != nil {
		sampling = options.SamplingInterval
		exactMatch = options.ExactMatch
	}
	pixels := extractColors(image, int(sampling))

	extractor := newDominantColorExtractor(number, exactMatch)
	return extractor.extractDominantColors(pixels)
}

func extractColors(img image.Image, sampling int) []Color {
	if sampling == 0 {
		longer := max(img.Bounds().Dx(), img.Bounds().Dy())
		sampling = longer / 64
	}
	w := img.Bounds().Dx() / sampling
	h := img.Bounds().Dy() / sampling

	colors := make([]Color, w*h)
	index := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := NewColourFromColor(img.At(x*sampling, y*sampling))
			colors[index] = c
			index++
		}
	}
	return colors
}
