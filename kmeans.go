package dcolors

import (
	"math/rand"
	"slices"
)

type InitialSelectionType int

const (
	// UniformSelection is used to select initial pixels from the image with a uniform spread over the image
	UniformSelection InitialSelectionType = 0
	// RandomSelection is used to select random pixels for the initial seed
	RandomSelection = 1
)

type colorExtractor struct {
	selectionType InitialSelectionType
	numCentroids  int
	pixels        []Color
	exactMatch    bool
}

func newDominantColorExtractor(numCentroids int, selectionType InitialSelectionType, match bool) *colorExtractor {
	return &colorExtractor{numCentroids: numCentroids, exactMatch: match, selectionType: selectionType}
}

func (c *colorExtractor) extractDominantColors(pixels []Color) []Color {
	c.pixels = pixels

	return c.runKMeans()
}

// runKmeans executes a K-means algorithm
// this implementation strives to minimize memory allocations
func (c *colorExtractor) runKMeans() []Color {
	// initialize centroids
	var centroids []Color
	if c.selectionType == UniformSelection {
		centroids = c.spreadCentroids()
	} else {
		centroids = c.randomCentroids()
	}
	newCentroids := make([]Color, len(centroids))
	copy(newCentroids, centroids)

	// prepare cluster color aggregator
	// cluster color sum contains sum of RGB values in the cluster
	clusterColorSum := make([][3]uint32, len(centroids))
	// clusterCounts contains the number of pixels assigned to this cluster
	clusterCounts := make([]uint32, len(centroids))
	// next definitions are only used if we're looking for exact color from image
	var pixelIndices []int
	var closestPixelIdx []int
	var closestDistance []float64
	if c.exactMatch {
		// we will need to keep track of the cluster to which each color belongs
		pixelIndices = make([]int, len(c.pixels))
		closestPixelIdx = make([]int, len(centroids))
		closestDistance = make([]float64, len(centroids))
	}

	// this functions clears the average color of each cluster and the number of pixels in it
	resetClusters := func() {
		for i, _ := range clusterColorSum {
			clusterColorSum[i][0] = 0
			clusterColorSum[i][1] = 0
			clusterColorSum[i][2] = 0
			clusterCounts[i] = 0
		}
	}
	// this function sums the RGB values of the color of the cluster and a provided color
	sumRgb := func(target *[3]uint32, c Color) {
		target[0] += c.rgb[0]
		target[1] += c.rgb[1]
		target[2] += c.rgb[2]
	}

	converged := false
	for !converged {
		for pxIndex, point := range c.pixels {
			distances := c.distanceToCentroids(point, centroids)
			clusterIndex := c.bestClusterIndex(distances)
			// add the rgb values of the color of the cluster and increment the pixel count for that cluster
			sumRgb(&clusterColorSum[clusterIndex], point)
			clusterCounts[clusterIndex]++
			if c.exactMatch {
				// store the index of the cluster for this pixel
				pixelIndices[pxIndex] = clusterIndex
			}
		}
		// calculate new centroid color for each cluster
		for i, cluster := range clusterColorSum {
			N := clusterCounts[i]
			newCentroids[i] = NewColorFromRgb(cluster[0]/N, cluster[1]/N, cluster[2]/N)
		}
		if c.exactMatch {
			// initialize selected pixel and distance to the calculated centroid
			for i := range newCentroids {
				closestDistance[i] = 99e100
				closestPixelIdx[i] = -1
			}
			// for every pixel find the closest centroid. If current pixel is closer to the centroid
			// then the previous one, select this pixel as possible new centroid
			for i, clusterIndex := range pixelIndices {
				currentCentroid := newCentroids[clusterIndex]
				distance := currentCentroid.Distance(c.pixels[i])
				if distance < closestDistance[clusterIndex] {
					closestDistance[clusterIndex] = distance
					closestPixelIdx[clusterIndex] = i
				}
			}
			// replace calculate centroid color with a closest exact color from the image
			for i, closestColor := range closestPixelIdx {
				newCentroids[i] = c.pixels[closestColor]
			}
		}
		converged = slices.Equal(centroids, newCentroids)
		if !converged {
			// we need another iteration - clear cluster RGB sums and pixel counts
			resetClusters()
		}
		copy(centroids, newCentroids)
	}
	// sort cluster by the number of pixels assigned to each cluster, largest first
	centroidIndices := make([]int, len(centroids))
	for i := range centroidIndices {
		centroidIndices[i] = i
	}
	slices.SortFunc(centroidIndices, func(a, b int) int {
		return int(clusterCounts[b]) - int(clusterCounts[a])
	})

	results := make([]Color, len(centroids))
	for targetIdx, srcIdx := range centroidIndices {
		results[targetIdx] = centroids[srcIdx]
	}
	return results
}

// spreadCentroids samples the image to generate initial cluster centroids
func (c *colorExtractor) spreadCentroids() []Color {
	centroids := make([]Color, c.numCentroids)
	step := len(c.pixels) / c.numCentroids
	idx := 0
	for i := range centroids {
		centroids[i] = c.pixels[idx]
		idx += step
	}
	return centroids
}

// randomCentroids samples the image to generate initial cluster centroids
func (c *colorExtractor) randomCentroids() []Color {
	centroids := make([]Color, c.numCentroids)
	usedIndices := make(map[int]struct{})
	i := 0
	for {
		idx := rand.Intn(c.numCentroids)
		if _, present := usedIndices[idx]; present {
			continue
		}
		usedIndices[idx] = struct{}{}
		centroids[i] = c.pixels[idx]
		i++
		if i == len(centroids) {
			break
		}
	}
	return centroids
}

func (c *colorExtractor) distanceToCentroids(point Color, centroids []Color) []float64 {
	distances := make([]float64, c.numCentroids)
	for i, centroid := range centroids {
		distances[i] = centroid.Distance(point)
	}
	return distances
}

func (c *colorExtractor) bestClusterIndex(distances []float64) int {
	minDistance := 99e99
	minIndex := -1
	for i, d := range distances {
		if d < minDistance {
			minDistance = d
			minIndex = i
		}
	}
	return minIndex
}
