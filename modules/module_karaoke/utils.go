package module_karaoke

import (
	"sort"

	"github.com/partyhall/partyhall/services"
)

func getBestImage(images []services.SpotifyImage) (bestImage *services.SpotifyImage) {
	// Sort images by resolution (descending order)
	sort.Slice(images, func(i, j int) bool {
		areaI := images[i].Width * images[i].Height
		areaJ := images[j].Width * images[j].Height
		return areaI > areaJ
	})

	// Find the first image with size 300x300
	for _, img := range images {
		if img.Width == 300 && img.Height == 300 {
			return &img
		}
	}

	// If no 300x300 image, return the highest resolution image
	if len(images) > 0 {
		return &images[0]
	}

	// Return an empty image if the input array is empty
	return nil
}
