package constants

type LitImage struct {
	Name    string
	ImageID string
}

var knownImages []LitImage

func KnownImages() []LitImage {
	if (len(knownImages)) == 0 {
		knownImages = append(knownImages, LitImage{"Default", "sha256:2d6737a14c759c4d32b35187e6056430d5fc51592b1bd0d1b2713b85c853efbf"})
	}
	return knownImages
}

func DefaultImage() LitImage {
	for _, i := range KnownImages() {
		if i.Name == "Default" {
			return i
		}
	}
	return LitImage{}
}
