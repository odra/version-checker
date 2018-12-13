package github

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"github.com/integr8ly/version-checker/pkg/credentials"
	"github.com/integr8ly/version-checker/pkg/meta"
	"net/http"
	"os"
	"testing"
)

func TestNewGitHub(t *testing.T) {
	cases := []struct {
		Name     string
		Version  string
		Client   func(version string) *GitHub
		Validate func(t *testing.T, version string, client *GitHub)
	}{
		{
			Name:    "Should validate its beta version",
			Version: "beta",
			Client: func(version string) *GitHub {
				return NewGitHub(version)
			},
			Validate: func(t *testing.T, version string, client *GitHub) {
				clientVersion := client.Version
				expectedVersion := version
				if clientVersion.Ver != expectedVersion {
					t.Fatalf("Expected version %s but got %s", expectedVersion, clientVersion)
				}
			},
		},
		{
			Name:    "Should not validate its beta version",
			Version: "alpha",
			Client: func(version string) *GitHub {
				return NewGitHub(version)
			},
			Validate: func(t *testing.T, version string, client *GitHub) {
				clientVersion := client.Version
				if clientVersion.Ver == "beta" {
					t.Fatalf("Expected version alpha but got %s", clientVersion)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			client := tc.Client(tc.Version)
			tc.Validate(t, tc.Version, client)
		})
	}
}

func TestGitHub_Bootstrap(t *testing.T) {
	cases := []struct {
		Name     string
		Init     func()
		Clean    func()
		Client   func() *GitHub
		Validate func(t *testing.T, c *GitHub)
	}{
		{
			Name:  "Should fail to bootstrap",
			Init:  func() {},
			Clean: func() {},
			Client: func() *GitHub {
				return NewGitHub("0.0.1")
			},
			Validate: func(t *testing.T, c *GitHub) {
				err := c.Bootstrap()
				if err == nil {
					t.Fatalf("Client bootstrap should have failed")
				}
			},
		},
		{
			Name: "Should bootstrap client",
			Init: func() {
				os.Setenv("GITHUB_APP_TOKEN", "github")
				credentials.Bootstrap()
			},
			Clean: func() {
				os.Unsetenv("GITHUB_APP_TOKEN")
			},
			Client: func() *GitHub {
				return NewGitHub("0.0.1")
			},
			Validate: func(t *testing.T, c *GitHub) {
				err := c.Bootstrap()
				if err != nil {
					t.Fatalf("Failed to bootstrap client: %v", err)
				}

				if c.Client == nil {
					t.Fatal("Githug rest client is nil")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Init()
			c := tc.Client()
			tc.Validate(t, c)
			tc.Clean()
		})
	}
}

func TestGitHub_IsKind(t *testing.T) {
	cases := []struct {
		Name     string
		Client   func() *GitHub
		Kind     meta.ReleaseKind
		Validate func(t *testing.T, kind meta.ReleaseKind, client *GitHub)
	}{
		{
			Name: "Shoud validate correct kind",
			Kind: GHKind,
			Client: func() *GitHub {
				return NewGitHub("beta")
			},
			Validate: func(t *testing.T, kind meta.ReleaseKind, client *GitHub) {
				if !client.IsKind(kind) {
					t.Fatalf("Could not validate kind: %s", kind)
				}
			},
		},
		{
			Name: "Shoud not validate kind",
			Kind: "gitlab",
			Client: func() *GitHub {
				return NewGitHub("beta")
			},
			Validate: func(t *testing.T, kind meta.ReleaseKind, client *GitHub) {
				if client.IsKind(kind) {
					t.Fatalf("Should not validate kind: %s", kind)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c := tc.Client()
			tc.Validate(t, tc.Kind, c)
		})
	}
}

func TestGitHub_LatestVersion(t *testing.T) {
	cases := []struct {
		Name     string
		Version  string
		Client   func() (meta.ReleaseSpec, func())
		Validate func(t *testing.T, version string, c meta.ReleaseSpec)
	}{
		{
			Name:    "Should retrieve latest version",
			Version: "0.0.2",
			Client: func() (meta.ReleaseSpec, func()) {
				c := NewGitHub("0.0.1")
				c.Org = "o"
				c.Repo = "r"

				fakeClient, mux, _, teardown := setupFakeClient()

				mux.HandleFunc("/repos/o/r/releases/latest", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					tag := new(string)
					*tag = "0.0.2"
					ghRelease := github.RepositoryRelease{
						TagName: tag,
					}

					b, err := json.Marshal(ghRelease)
					if err != nil {
						t.Fatalf("Error while serializing json data: %v", err)
					}

					w.Write(b)
				})

				c.Client = &client{
					restClient: fakeClient,
				}

				return c, teardown
			},
			Validate: func(t *testing.T, version string, c meta.ReleaseSpec) {
				ver, err := c.LatestVersion()
				if err != nil {
					t.Fatalf("Failed to retrieve version: %v", err)
				}

				if ver.Ver != version {
					t.Fatalf("Expected version %s but got %s", version, ver.Ver)
				}
			},
		},
		{
			Name:    "Should not retrieve latest version",
			Version: "0.0.2",
			Client: func() (meta.ReleaseSpec, func()) {
				c := NewGitHub("0.0.1")
				c.Org = "o"
				c.Repo = "r"

				fakeClient, mux, _, teardown := setupFakeClient()

				mux.HandleFunc("/repos/o/r/releases", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					tag := new(string)
					*tag = "0.0.2"
					ghRelease := github.RepositoryRelease{
						TagName: tag,
					}

					b, err := json.Marshal(ghRelease)
					if err != nil {
						t.Fatalf("Error while serializing json data: %v", err)
					}

					w.Write(b)
				})

				c.Client = &client{
					restClient: fakeClient,
				}

				return c, teardown
			},
			Validate: func(t *testing.T, version string, c meta.ReleaseSpec) {
				_, err := c.LatestVersion()
				if err == nil {
					t.Fatal("Should fail to retrieve latest version")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c, teardown := tc.Client()
			tc.Validate(t, tc.Version, c)
			teardown()
		})
	}
}

func TestGitHub_HasNewVersion(t *testing.T) {
	cases := []struct {
		Name     string
		Version  string
		Client   func(ver string) (meta.ReleaseSpec, func())
		Validate func(t *testing.T, c meta.ReleaseSpec)
	}{
		{
			Name:    "Should have new version",
			Version: "0.0.1",
			Client: func(ver string) (meta.ReleaseSpec, func()) {
				c := NewGitHub(ver)
				c.Org = "o"
				c.Repo = "r"

				fakeClient, mux, _, teardown := setupFakeClient()

				mux.HandleFunc("/repos/o/r/releases/latest", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					tag := new(string)
					*tag = "0.0.2"
					ghRelease := github.RepositoryRelease{
						TagName: tag,
					}

					b, err := json.Marshal(ghRelease)
					if err != nil {
						t.Fatalf("Error while serializing json data: %v", err)
					}

					w.Write(b)
				})

				c.Client = &client{
					restClient: fakeClient,
				}

				return c, teardown
			},
			Validate: func(t *testing.T, c meta.ReleaseSpec) {
				ver, ok, err := c.HasNewVersion()
				if err != nil {
					t.Fatalf("Error while checking for new version: %v", err)
				}

				if !ok {
					t.Fatalf("New version is not available")
				}

				if ver.Ver != "0.0.2" {
					t.Fatalf("New version does not match: %v", ver)
				}
			},
		},
		{
			Name:    "Should not have new version",
			Version: "0.0.1",
			Client: func(ver string) (meta.ReleaseSpec, func()) {
				c := NewGitHub(ver)
				c.Org = "o"
				c.Repo = "r"

				fakeClient, mux, _, teardown := setupFakeClient()

				mux.HandleFunc("/repos/o/r/releases/latest", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					tag := new(string)
					*tag = "0.0.1"
					ghRelease := github.RepositoryRelease{
						TagName: tag,
					}

					b, err := json.Marshal(ghRelease)
					if err != nil {
						t.Fatalf("Error while serializing json data: %v", err)
					}

					w.Write(b)
				})

				c.Client = &client{
					restClient: fakeClient,
				}

				return c, teardown
			},
			Validate: func(t *testing.T, c meta.ReleaseSpec) {
				ver, ok, err := c.HasNewVersion()
				if err != nil {
					t.Fatalf("Error while checking for new version: %v", err)
				}

				if ok {
					t.Fatalf("Should not have a new version available")
				}

				if ver.Ver != "0.0.1" {
					t.Fatalf("New version does not match: %v", ver)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c, teardown := tc.Client(tc.Version)
			tc.Validate(t, c)
			teardown()
		})
	}
}
