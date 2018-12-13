package github

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const baseURLPath = "/api-v3"

func setupFakeClient() (client *github.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))

	server := httptest.NewServer(apiHandler)

	fakeClient := github.NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	fakeClient.BaseURL = url
	fakeClient.UploadURL = url

	return fakeClient, mux, server.URL, server.Close
}

func TestGHClient(t *testing.T) {
	cases := []struct {
		Name            string
		Release         string
		ExpectedRelease string
		Client          func(t *testing.T, release string) (*client, func())
		Validate        func(t *testing.T, release string, c *client)
	}{
		{
			Name:            "Should retrieve correct tag",
			Release:         "latest-release",
			ExpectedRelease: "latest-release",
			Client: func(t *testing.T, release string) (*client, func()) {
				fakeClient, mux, _, teardown := setupFakeClient()

				mux.HandleFunc("/repos/o/r/releases/latest", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					tag := new(string)
					*tag = release
					ghRelease := github.RepositoryRelease{
						TagName: tag,
					}

					b, err := json.Marshal(ghRelease)
					if err != nil {
						t.Fatalf("Error while serializing json data: %v", err)
					}

					w.Write(b)
				})

				c := &client{
					restClient: fakeClient,
				}

				return c, teardown
			},
			Validate: func(t *testing.T, release string, c *client) {
				r, e := c.GetLatestRelease("o", "r")
				if e != nil {
					t.Fatalf("Failed to get latest release: %v", e)
				}

				actual := string(*r.TagName)
				expected := string(release)

				if strings.Compare(actual, expected) != 0 {
					t.Fatalf("Expected tag %s but got %s", release, *r.TagName)
				}
			},
		},
		{
			Name:            "Should retrieve incorrect tag",
			Release:         "latest-release",
			ExpectedRelease: "another-tag",
			Client: func(t *testing.T, release string) (*client, func()) {
				fakeClient, mux, _, teardown := setupFakeClient()

				mux.HandleFunc("/repos/o/r/releases/latest", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					tag := new(string)
					*tag = release
					ghRelease := github.RepositoryRelease{
						TagName: tag,
					}

					b, err := json.Marshal(ghRelease)
					if err != nil {
						t.Fatalf("Error while serializing json data: %v", err)
					}

					w.Write(b)
				})

				c := &client{
					restClient: fakeClient,
				}

				return c, teardown
			},
			Validate: func(t *testing.T, release string, c *client) {
				r, e := c.GetLatestRelease("o", "r")
				if e != nil {
					t.Fatalf("Failed to get latest release: %v", e)
				}

				actual := string(*r.TagName)
				expected := string(release)

				if strings.Compare(actual, expected) == 0 {
					t.Fatalf("Expected tag %s but got %s", release, *r.TagName)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c, teardown := tc.Client(t, tc.ExpectedRelease)
			tc.Validate(t, tc.Release, c)
			teardown()
		})
	}
}
