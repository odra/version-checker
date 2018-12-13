package github

import (
	"errors"
	"github.com/integr8ly/version-checker/pkg/credentials"
	"github.com/integr8ly/version-checker/pkg/meta"
)

const GHKind meta.ReleaseKind = "github"

type GitHub struct {
	Kind    meta.ReleaseKind
	Version meta.Version
	Org     string
	Repo    string
	Client  GHClientSpec
}

func NewGitHub(version string) *GitHub {
	return &GitHub{
		Kind: GHKind,
		Version: meta.Version{
			Ver: version,
		},
	}
}

func (gh *GitHub) Bootstrap() error {
	cred := credentials.Get()
	if cred == nil || cred.GitHub.Token == "" {
		return errors.New("github token is not set")
	}

	gh.Client = GHClient(cred.GitHub.Token)

	return nil
}

func (gh *GitHub) IsKind(kind meta.ReleaseKind) bool {
	return gh.Kind == kind
}

func (gh *GitHub) LatestVersion() (*meta.Version, error) {
	release, err := gh.Client.GetLatestRelease(gh.Org, gh.Repo)
	if err != nil {
		return nil, err
	}

	return &meta.Version{
		Ver: *release.TagName,
	}, nil
}

func (gh *GitHub) HasNewVersion() (*meta.Version, bool, error) {
	latest, err := gh.LatestVersion()
	if err != nil {
		return nil, false, err
	}

	if latest.Ver == gh.Version.Ver {
		return latest, false, nil
	}

	return latest, true, nil
}
