package constants

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type LitImage struct {
	Name    string
	ImageID string
}

var knownImages []LitImage

func KnownImages() []LitImage {
	return knownImages
}

func InitImages(cli *client.Client) error {
	res, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		return err
	}

	for _, r := range res {
		for _, tag := range r.RepoTags {
			if tag[:4] == "lit:" {
				name := "Default"
				if tag != "lit:latest" {
					name = tag[4:]
				}
				knownImages = append(knownImages, LitImage{ImageID: r.ID, Name: name})
			}
		}
	}
	return nil
}

func DefaultImage() LitImage {
	for _, i := range KnownImages() {
		if i.Name == "Default" {
			return i
		}
	}
	return LitImage{}
}
