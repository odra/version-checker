package credentials

import (
	"github.com/integr8ly/version-checker/pkg/meta"
	"os"
	"testing"
)

func cleanAll() {
	os.Unsetenv("GITHUB_APP_TOKEN")
	os.Unsetenv("DOCKER_USERNAME")
	os.Unsetenv("DOCKER_PASSWORD")
	os.Unsetenv("TRELLO_APP_KEY")
	os.Unsetenv("TRELLO_APP_TOKEN")
	Reset()
}

func TestBootstrap(t *testing.T) {
	cases := []struct {
		Name      string
		Bootstrap func()
		Validate  func(t *testing.T, c *meta.Credentials)
	}{
		{
			Name: "Should bootstrap all vars",
			Bootstrap: func() {
				os.Setenv("GITHUB_APP_TOKEN", "github")
				os.Setenv("DOCKER_USERNAME", "docker")
				os.Setenv("DOCKER_PASSWORD", "docker")
				os.Setenv("TRELLO_APP_KEY", "trello")
				os.Setenv("TRELLO_APP_TOKEN", "trello")
				Bootstrap()
			},
			Validate: func(t *testing.T, c *meta.Credentials) {
				if c.GitHub == nil || c.GitHub.Token != "github" {
					t.Fatalf("Failed to retrieve github token")
				}

				if c.Docker == nil || c.Docker.Username != "docker" || c.Docker.Password != "docker" {
					t.Fatalf("Failed to retrieve docker credentials")
				}

				if c.Trello == nil {
					t.Fatalf("Trello credentials not set")
				}

				if c.Trello.AppKey != "trello" || c.Trello.Token != "trello" {
					t.Fatalf("Failed to retrive trello credentials")
				}
			},
		},
		{
			Name: "Should bootstrap github only",
			Bootstrap: func() {
				os.Setenv("GITHUB_APP_TOKEN", "github")
				Bootstrap()
			},
			Validate: func(t *testing.T, c *meta.Credentials) {
				if c.GitHub == nil || c.GitHub.Token != "github" {
					t.Fatalf("Failed to retrieve github token")
				}

				if c.Docker == nil || c.Docker.Username != "" || c.Docker.Password != "" {
					t.Fatalf("Docker credentials should not be set")
				}

				if c.Trello == nil {
					t.Fatalf("Trello credentials not set")
				}

				if c.Trello.AppKey != "" || c.Trello.Token != "" {
					t.Fatalf("Trello credentials should be empty")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Bootstrap()
			tc.Validate(t, Get())
			cleanAll()
		})
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		Name      string
		Bootstrap func()
		Validate  func(t *testing.T, c *meta.Credentials)
	}{
		{
			Name: "Should return a nil client",
			Bootstrap: func() {

			},
			Validate: func(t *testing.T, c *meta.Credentials) {
				if c != nil {
					t.Fatalf("Credentials should be nil but got: %v", c)
				}
			},
		},
		{
			Name: "Shoud retrieve a client",
			Bootstrap: func() {
				os.Setenv("GITHUB_APP_TOKEN", "github")
				os.Setenv("DOCKER_USERNAME", "docker")
				os.Setenv("DOCKER_PASSWORD", "docker")
				os.Setenv("TRELLO_APP_KEY", "trello")
				os.Setenv("TRELLO_APP_TOKEN", "trello")
				Bootstrap()
			},
			Validate: func(t *testing.T, c *meta.Credentials) {
				if c == nil {
					t.Fatalf("Credentials should not be nil")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Bootstrap()
			tc.Validate(t, Get())
			cleanAll()
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		Name      string
		Bootstrap func()
		Cred      func() *meta.Credentials
		Validate  func(t *testing.T, c meta.Credentials)
	}{
		{
			Name: "Update should fail due to nil credentials",
			Bootstrap: func() {

			},
			Cred: func() *meta.Credentials {
				return &meta.Credentials{
					GitHub: &meta.GitHubCredential{
						Token: "github",
					},
				}
			},
			Validate: func(t *testing.T, c meta.Credentials) {
				err := Update(c)
				if err == nil {
					t.Fatalf("Should get error but credential is not nil: %v", c)
				}
			},
		},
		{
			Name: "Should update github token",
			Bootstrap: func() {
				os.Setenv("GITHUB_APP_TOKEN", "github")
				os.Setenv("DOCKER_USERNAME", "docker")
				os.Setenv("DOCKER_PASSWORD", "docker")
				os.Setenv("TRELLO_APP_KEY", "trello")
				os.Setenv("TRELLO_APP_TOKEN", "trello")
				Bootstrap()
			},
			Cred: func() *meta.Credentials {
				return &meta.Credentials{
					GitHub: &meta.GitHubCredential{
						Token: "anoter_github_token",
					},
				}
			},
			Validate: func(t *testing.T, c meta.Credentials) {
				err := Update(c)
				if err != nil {
					t.Fatalf("Should get error but credential is not nil: %v", c)
				}

				if c.GitHub == nil || c.GitHub.Token != "anoter_github_token" {
					t.Fatalf("Expected github token \"anoter_github_token\" but got \"%s\"", c.GitHub.Token)
				}

			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Bootstrap()
			c := tc.Cred()
			tc.Validate(t, *c)
			cleanAll()
		})
	}
}

func TestReset(t *testing.T) {
	cases := []struct {
		Name      string
		Bootstrap func()
		Validate  func(t *testing.T, c *meta.Credentials)
	}{
		{
			Name: "Should reset credentials",
			Bootstrap: func() {
				os.Setenv("GITHUB_APP_TOKEN", "github")
				os.Setenv("DOCKER_USERNAME", "docker")
				os.Setenv("DOCKER_PASSWORD", "docker")
				os.Setenv("TRELLO_APP_KEY", "trello")
				os.Setenv("TRELLO_APP_TOKEN", "trello")
				Bootstrap()
			},
			Validate: func(t *testing.T, c *meta.Credentials) {
				if c != nil {
					t.Fatalf("Credentials should be nil but got %v", c)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Bootstrap()
			Reset()
			tc.Validate(t, Get())
			cleanAll()
		})
	}
}
