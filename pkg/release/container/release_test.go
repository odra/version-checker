package container

import (
	"github.com/integr8ly/version-checker/pkg/credentials"
	"github.com/integr8ly/version-checker/pkg/meta"
	"os"
	"testing"
)

func TestNewContainer(t *testing.T) {
	cases := []struct {
		Name      string
		Container func() *Container
		Validate  func(t *testing.T, c *Container)
	}{
		{
			Name: "Should create a new container release",
			Container: func() *Container {
				return NewContainer("0.0.1")
			},
			Validate: func(t *testing.T, c *Container) {
				if c.Kind != "container-registry" {
					t.Fatalf("Kind does not match: %s", c.Kind)
				}

				if c.Version.Ver != "0.0.1" {
					t.Fatalf("Version does not match: %s", c.Version.Ver)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c := tc.Container()
			tc.Validate(t, c)
		})
	}
}

func TestContainer_Bootstrap(t *testing.T) {
	cases := []struct {
		Name      string
		Init      func()
		Container func() *Container
		Clean     func()
		Validate  func(t *testing.T, c *Container)
	}{
		{
			Name: "Should bootstrap container",
			Container: func() *Container {
				return NewContainer("0.0.1")
			},
			Init: func() {
				os.Setenv("DOCKER_USERNAME", "docker")
				os.Setenv("DOCKER_PASSWORD", "docker")
				credentials.Bootstrap()
			},
			Clean: func() {
				os.Unsetenv("DOCKER_USERNAME")
				os.Unsetenv("DOCKER_PASSWORD")
				credentials.Reset()
			},
			Validate: func(t *testing.T, c *Container) {
				err := c.Bootstrap()
				if err != nil {
					t.Fatalf("Could not bootrap container: %v", err)
				}
			},
		},
		{
			Name: "Should not bootstrap container",
			Container: func() *Container {
				return NewContainer("0.0.1")
			},
			Init: func() {

			},
			Clean: func() {

			},
			Validate: func(t *testing.T, c *Container) {
				err := c.Bootstrap()
				if err == nil {
					t.Fatal("Container bootstrap should fail")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c := tc.Container()
			tc.Init()
			tc.Validate(t, c)
			tc.Clean()
		})
	}
}

func TestContainer_IsKind(t *testing.T) {
	cases := []struct {
		Name      string
		Container func() *Container
		Validate  func(t *testing.T, c *Container)
	}{
		{
			Name: "Should validate kind",
			Container: func() *Container {
				return NewContainer("0.0.1")
			},
			Validate: func(t *testing.T, c *Container) {
				if !c.IsKind(ContainerKind) {
					t.Fatalf("Kind does not match: %s", c.Kind)
				}
			},
		},
		{
			Name: "Should not validate kind",
			Container: func() *Container {
				return NewContainer("0.0.1")
			},
			Validate: func(t *testing.T, c *Container) {
				var kind meta.ReleaseKind = "container-kind"
				if c.IsKind(kind) {
					t.Fatalf("Kind should not match: %s", c.Kind)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c := tc.Container()
			tc.Validate(t, c)
		})
	}
}

func TestContainer_LatestVersion(t *testing.T) {
	cases := []struct {
		Name      string
		Version   string
		Container func() *Container
		Validate  func(t *testing.T, c *Container, version string)
	}{
		{
			Name:    "Should validate version",
			Version: "0.0.1",
			Container: func() *Container {
				return NewContainer("0.0.1")
			},
			Validate: func(t *testing.T, c *Container, version string) {
				v, e := c.LatestVersion()
				if e != nil {
					t.Fatalf("failed to get latest version: %v", e)
				}

				if v.Ver != version {
					t.Fatalf("Version does not match: %s", v.Ver)
				}
			},
		},
		{
			Name:    "Should not validate version",
			Version: "0.0.2",
			Container: func() *Container {
				return NewContainer("0.0.1")
			},
			Validate: func(t *testing.T, c *Container, version string) {
				v, e := c.LatestVersion()
				if e != nil {
					t.Fatalf("failed to get latest version: %v", e)
				}

				if v.Ver == version {
					t.Fatalf("Version should not match: %s", v.Ver)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c := tc.Container()
			tc.Validate(t, c, tc.Version)
		})
	}
}

func TestContainer_HasNewVersion(t *testing.T) {
	cases := []struct {
		Name      string
		Version   string
		Container func() *Container
		Validate  func(t *testing.T, c *Container, version string)
	}{
		{
			Name:    "Should identify new version",
			Version: "0.0.2",
			Container: func() *Container {
				return NewContainer("0.0.1")
			},
			Validate: func(t *testing.T, c *Container, version string) {
				_, _, err := c.HasNewVersion()
				if err != nil {
					t.Fatalf("Could not verify new version: %v", err)
				}

				//if !has {
				//	t.Fatalf("Could not identify a new version: %s", ver.Ver)
				//}
			},
		},
		{
			Name:    "Should not identify new version",
			Version: "0.0.1",
			Container: func() *Container {
				return NewContainer("0.0.1")
			},
			Validate: func(t *testing.T, c *Container, version string) {
				ver, has, err := c.HasNewVersion()
				if err != nil {
					t.Fatalf("Could not verify new version: %v", err)
				}

				if has {
					t.Fatalf("Should not identify a new version: %s", ver.Ver)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c := tc.Container()
			tc.Validate(t, c, tc.Version)
		})
	}
}
