package container

import (
	"errors"
	"github.com/integr8ly/version-checker/pkg/credentials"
	"github.com/integr8ly/version-checker/pkg/meta"
	"github.com/sirupsen/logrus"
	"strings"
)

const ContainerKind meta.ReleaseKind = "container-registry"

type Container struct {
	Kind    meta.ReleaseKind
	Version meta.Version
	Org     string
	Name    string
	Client  ClientSpec
}

func NewContainer(version string) *Container {
	return &Container{
		Kind: ContainerKind,
		Version: meta.Version{
			Ver: version,
		},
	}
}

func (c *Container) Bootstrap() error {
	cred := credentials.Get()
	if cred == nil || cred.Docker.Username == "" || cred.Docker.Password == "" {
		return errors.New("container token is not set")
	}

	return nil
}

func (c *Container) IsKind(kind meta.ReleaseKind) bool {
	return c.Kind == kind
}

func (c *Container) LatestVersion() (*meta.Version, error) {
	return &meta.Version{
		Ver: "0.0.1",
	}, nil
}

func (c *Container) HasNewVersion() (*meta.Version, bool, error) {
	latest, err := c.LatestVersion()
	if err != nil {
		return nil, false, err
	}

	if strings.Compare(latest.Ver, c.Version.Ver) == 0 {
		logrus.Infof("|||%s:%s|||", latest.Ver, c.Version.Ver)
		return latest, false, nil
	}

	return latest, true, nil
}
