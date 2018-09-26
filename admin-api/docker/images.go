package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetImage(cli *client.Client, name string) (string, error) {
	res, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		return "", err
	}

	for _, r := range res {
		for _, tag := range r.RepoTags {
			if tag == fmt.Sprintf("%s:latest", name) {
				return r.ID, nil
			}
		}
	}
	return "", fmt.Errorf("Image not found: %s", name)
}
