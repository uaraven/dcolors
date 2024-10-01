package dcolors

import "image"

// Options contains configuration parameters for color extractor
//
// SamplingInterval - selects how many pixels from the original image is used for color extraction.
// It is rarely needed to process every pixel in the image, so sampling only some pixels may drastically improve
// performance. SamplingInterval defines how many pixels to skip between the chosen one. It essentially scales the image
// down by SamplingInterval times. If '0' is provided then such SamplingInterval is selected that image is scaled to
// 64 pixels by longest side
//
// InitialSelection - configures the mode of selecting the initial pixels. Setting this parameter to UniformSelection
// will sample the initial colors uniformly from the image. Setting RandomSelection will choose initial pixels at random.
// Performance and results might differ between these two options, so test which one is better suited for your needs.
// Default value is UniformSelection
//
// ExactMatch - if set to true then resulting dominant colors will be selected from the exact colors in the image.
// Otherwise, the resulting colors will represent the average colors, but not necessary exact colors of the image.
// Selecting exact colors might be faster in some cases.
// Default value is false.
type Options struct {
	SamplingInterval uint
	InitialSelection InitialSelectionType
	ExactMatch       bool
}

// ExtractDominantColors will extract the most prominent colors from the image.
//
// image - the image.Image containing the pixel data of the picture to process
//
// number - the number of colors to extract.
//
// options - pointer to the Options structure containing tuning parameters for the algorithm
//
// Returns the slice of Color elements. The colors are sorted in the order of the number of amount of pixels in the image
// corresponding to the returned color. This way the most prominent colors are at the lowest indices of the returned slice.
//
// For example, if the 80-pixel image contains 20 red pixels and 60 blue pixels then
// the result will be {blue, red}
func ExtractDominantColors(image image.Image, number int, options *Options) []Color {
	var sampling uint
	exactMatch := false
	initialSelection := UniformSelection
	if options != nil {
		sampling = options.SamplingInterval
		exactMatch = options.ExactMatch
		initialSelection = options.InitialSelection
	}
	pixels := extractColors(image, int(sampling))

	extractor := newDominantColorExtractor(number, initialSelection, exactMatch)
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
			c := NewColorFromColor(img.At(x*sampling, y*sampling))
			colors[index] = c
			index++
		}
	}
	return colors
}
