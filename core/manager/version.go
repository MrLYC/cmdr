package manager

import (
	"fmt"
	"strconv"

	"github.com/asdine/storm/v3/q"
	ver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

func normalizeVersion(version string) (string, error) {
	if version == "" {
		return "", nil
	}

	semver, err := ver.NewVersion(version)
	if err != nil {
		return "", errors.Wrapf(err, "invalid version %s", version)
	}

	return semver.String(), nil
}

func displayVersion(version string) string {
	semver, err := ver.NewVersion(version)
	if err != nil {
		return version
	}

	segments := semver.Segments()
	if segments[1] == 0 && segments[2] == 0 {
		return strconv.Itoa(segments[0])
	}
	if segments[2] == 0 {
		return fmt.Sprintf("%d.%d", segments[0], segments[1])
	}
	return fmt.Sprintf("%d.%d.%d", segments[0], segments[1], segments[2])
}

func versionMatches(candidate, wanted string) (bool, error) {
	if candidate == wanted || displayVersion(candidate) == wanted {
		return true, nil
	}

	wantedVersion, err := ver.NewVersion(wanted)
	if err != nil {
		return false, errors.Wrapf(err, "invalid version %s", wanted)
	}

	candidateVersion, err := ver.NewVersion(candidate)
	if err != nil {
		return false, nil
	}

	return candidateVersion.Equal(wantedVersion), nil
}

func queryMatchVersion(version string) (q.Matcher, error) {
	normalized, err := normalizeVersion(version)
	if err != nil {
		return nil, err
	}

	return q.Or(
		q.Eq("Version", version),
		q.Eq("Version", normalized),
	), nil
}
