package container

import (
	"errors"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/heroku/docker-registry-client/registry"
	"reflect"
)

var containerClient *client

type client struct {
	hub HubSpec
}

type ClientSpec interface {
	GetTags(repo string) ([]string, error)
	GetLatestTag(tags []string) (string, error)
	GetManifest(repo string, tag string) (*schema1.SignedManifest, error)
	CompareTags(repo string, newTag string, currentTag string) (int, error)
}

type HubSpec interface {
	Tags(repository string) (tags []string, err error)
	Manifest(repository, reference string) (*schema1.SignedManifest, error)
}

func Build(url string, username string, password string) (*client, error) {
	if containerClient != nil {
		return containerClient, nil
	}

	hub, err := registry.New(url, username, password)
	if err != nil {
		return nil, err
	}

	containerClient = &client{
		hub: hub,
	}

	return containerClient, nil
}

func (c *client) GetTags(repo string) ([]string, error) {
	tags, err := c.hub.Tags(repo)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (c *client) GetLatestTag(tags []string) (string, error) {
	if len(tags) == 0 {
		return "", errors.New("tag list is empty")
	}

	return tags[0], nil
}

func (c *client) CompareTags(repo string, newTag string, currentTag string) (int, error) {
	newManifest, err := c.GetManifest(repo, newTag)
	if err != nil {
		return -1, err
	}

	currentManifest, err := c.GetManifest(repo, currentTag)
	if err != nil {
		return -1, err
	}

	if reflect.DeepEqual(newManifest, currentManifest) {
		return 0, nil
	}

	return 1, nil
}

func (c *client) GetManifest(repo string, tag string) (*schema1.SignedManifest, error) {
	manifest, err := c.hub.Manifest(repo, tag)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}
