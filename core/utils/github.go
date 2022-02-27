package utils

import (
	"context"

	"github.com/google/go-github/v39/github"

	"github.com/mrlyc/cmdr/core"
)

func GetCMDRRelease(ctx context.Context, tag string) (release *github.RepositoryRelease, err error) {
	client := github.NewClient(nil)
	if tag == "latest" {
		release, _, err = client.Repositories.GetLatestRelease(ctx, core.Author, core.Name)
	} else {
		release, _, err = client.Repositories.GetReleaseByTag(ctx, core.Author, core.Name, tag)
	}

	return
}
