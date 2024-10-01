package dcolors

import (
	"fmt"
	"testing"
)

func BenchmarkExtractor(b *testing.B) {
	var img = loadImage("test-data/ir1.jpg")

	colors := extractColors(img, 0)

	type testSpec struct {
		title     string
		colors    int
		exact     bool
		selection InitialSelectionType
	}
	testData := []testSpec{
		{"Default sampling,4 colors,non-exact", 4, false, UniformSelection},
		{"Default sampling,4 colors,exact", 4, true, UniformSelection},
		{"Default sampling,4 colors,non-exact,random", 4, false, RandomSelection},
		{"Default sampling,4 colors,exact,random", 4, true, RandomSelection},
		{"Default sampling,8 colors,non-exact", 8, false, UniformSelection},
		{"Default sampling,8 colors,exact", 8, true, UniformSelection},
		{"Default sampling,8 colors,non-exact,random", 8, false, RandomSelection},
		{"Default sampling,8 colors,exact,random", 8, true, RandomSelection},
	}
	for _, test := range testData {
		extractor := newDominantColorExtractor(test.colors, test.selection, test.exact)
		b.Run(test.title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				extractor.extractDominantColors(colors)
			}
		})
	}

}

func BenchmarkImageRgbReader(b *testing.B) {
	var img = loadImage("test-data/ir1.jpg")

	for sampling := 0; sampling < 100; sampling += 10 {
		b.Run(fmt.Sprintf("Extract colors. Sampling: %d", sampling), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				extractColors(img, sampling)
			}
		})
	}
}
