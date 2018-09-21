package constants

type LitImage struct {
	Name    string
	ImageID string
}

var knownImages []LitImage

func KnownImages() []LitImage {
	if (len(knownImages)) == 0 {
		//knownImages = append(knownImages, LitImage{"Default", "sha256:eda4b2e362fe90ce8da6fd3e8202b316b84013e2176b91989129257d60d05ac5"})
	}
	return knownImages
}

func SetDefaultImage(id string) {
	knownImages = append(knownImages, LitImage{"Default", id})
}

func DefaultImage() LitImage {
	for _, i := range KnownImages() {
		if i.Name == "Default" {
			return i
		}
	}
	return LitImage{}
}
