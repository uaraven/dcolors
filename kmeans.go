package dominant_colors

import (
	"reflect"
	"slices"
)

type Cluster []Color

func (c Cluster) mean() Color {
	if len(c) == 0 {
		return Color{}
	}
	var r, g, b uint32
	for _, c := range c {
		r = r + c.rgb[0]
		g = g + c.rgb[1]
		b = b + c.rgb[2]
	}
	return NewColorFromRgb(r/uint32(len(c)), g/uint32(len(c)), b/uint32(len(c)))
}

func (c Cluster) closestActual() Color {
	avg := c.mean()
	closest := c[0]
	closestDistance := 1e100
	for _, px := range c {
		pxDist := px.Distance(avg)
		if pxDist < closestDistance {
			closest = px
			closestDistance = pxDist
		}
	}
	return closest
}

type colorExtractor struct {
	numCentroids int
	pixels       []Color
	clusters     []Cluster
	exactMatch   bool
}

func newDominantColourExtractor(numCentroids int, match bool) *colorExtractor {
	return &colorExtractor{numCentroids: numCentroids, exactMatch: match}
}

func (c *colorExtractor) extractDominantColours(pixels []Color) []Color {
	c.pixels = pixels

	return c.runKMeans()
}

func (c *colorExtractor) runKMeans() []Color {
	centroids := c.randomCentroids()
	var clusters []Cluster
	converged := false
	for !converged {
		clusters = c.clearClusters()

		for _, point := range c.pixels {
			distances := c.distanceToCentroids(point, centroids)
			clusterIndex := c.bestClusterIndex(distances)
			clusters[clusterIndex] = append(clusters[clusterIndex], point)
		}
		newCentroids := make([]Color, len(centroids))
		for i, cluster := range clusters {
			newCentroids[i] = cluster.mean()
		}
		converged = reflect.DeepEqual(centroids, newCentroids)
		centroids = newCentroids
	}
	centroidIndices := make([]int, len(centroids))
	for i := range centroidIndices {
		centroidIndices[i] = i
	}
	slices.SortFunc(centroidIndices, func(a, b int) int {
		return len(clusters[b]) - len(clusters[a])
	})
	results := make([]Color, len(centroids))
	for targetIdx, srcIdx := range centroidIndices {
		results[targetIdx] = centroids[srcIdx]
	}
	return results
}

// randomCentroids samples the image to generate initial cluster centroids
func (c *colorExtractor) randomCentroids() []Color {
	centroids := make([]Color, c.numCentroids)
	step := len(c.pixels) / c.numCentroids
	idx := 0
	for i := range centroids {
		centroids[i] = c.pixels[idx]
		idx += step
	}
	return centroids
}

func (c *colorExtractor) clearClusters() []Cluster {
	result := make([]Cluster, c.numCentroids)
	for i := range result {
		result[i] = []Color{}
	}
	return result
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
