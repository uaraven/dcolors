package dominant_colors

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"reflect"
	"testing"
)

type testData struct {
	image          string
	sampling       uint
	exactMatch     bool
	count          int
	expectedResult []Color
}

func TestExtractDominantColours(t *testing.T) {
	testData := []testData{
		{"test-data/cat.jpg", 0, false, 5,
			[]Color{
				NewColorFromRgb(28043, 21540, 24186),
				NewColorFromRgb(18939, 26594, 15474),
				NewColorFromRgb(39708, 38554, 29321),
				NewColorFromRgb(11587, 10290, 13257),
				NewColorFromRgb(1373, 7906, 16939),
			}},
		{"test-data/cat.jpg", 50, false, 5,
			[]Color{
				NewColorFromRgb(18680, 15231, 18523),
				NewColorFromRgb(29976, 35192, 23937),
				NewColorFromRgb(36771, 27766, 28650),
				NewColorFromRgb(13110, 18821, 7995),
				NewColorFromRgb(1439, 6828, 14026),
			}},
		{"test-data/cat.jpg", 0, true, 5,
			[]Color{
				NewColorFromRgb(28043, 21540, 24186),
				NewColorFromRgb(18939, 26594, 15474),
				NewColorFromRgb(39708, 38554, 29321),
				NewColorFromRgb(11587, 10290, 13257),
				NewColorFromRgb(1373, 7906, 16939),
			}},
	}

	for _, data := range testData {
		img := loadImage(data.image)
		result := ExtractDominantColours(img, data.count, &Options{
			SamplingInterval: data.sampling,
			ExactMatch:       data.exactMatch,
		})

		if !reflect.DeepEqual(result, data.expectedResult) {
			printFailure(result, data)
			t.Fail()
		}
	}

}

func printFailure(result []Color, data testData) {
	fmt.Printf("Test case image='%s', colours='%d', sampling='%d', exact match='%t' failed\n", data.image, data.count, data.sampling, data.exactMatch)
	fmt.Println("Expected results:")
	printColours(data.expectedResult)
	fmt.Println("Actual results:")
	printColours(result)
}

func printColours(data []Color) {
	for i, r := range data {
		// fmt.Print(r.Hex())
		fmt.Printf("NewColorFromRgb(%d,%d,%d)\n", r.rgb[0], r.rgb[1], r.rgb[2])
		if i == len(data)-1 {
			fmt.Println()
		} else {
			// fmt.Print(", ")
		}
	}
}

func loadImage(imagePath string) image.Image {
	f, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
