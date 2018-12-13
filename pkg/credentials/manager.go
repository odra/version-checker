package credentials

import (
	"errors"
	"github.com/integr8ly/version-checker/pkg/meta"
	"os"
)

var credentials *meta.Credentials

func Bootstrap() {
	credentials = &meta.Credentials{
		GitHub: &meta.GitHubCredential{
			Token: getEnv("GITHUB_APP_TOKEN"),
		},
		Docker: &meta.DockerCredential{
			Username: getEnv("DOCKER_USERNAME"),
			Password: getEnv("DOCKER_PASSWORD"),
		},
		Trello: &meta.TrelloCredential{
			AppKey: getEnv("TRELLO_APP_KEY"),
			Token:  getEnv("TRELLO_APP_TOKEN"),
		},
	}
}

func Get() *meta.Credentials {
	return credentials
}

func Update(c meta.Credentials) error {
	if credentials == nil {
		return errors.New("credential module needs to be boostraped first")
	}

	if c.GitHub != nil {
		credentials.GitHub.Token = c.GitHub.Token
	}

	if c.Docker != nil {
		credentials.Docker.Username = c.Docker.Username
		credentials.Docker.Password = c.Docker.Password
	}

	if c.Trello != nil {
		credentials.Trello.AppKey = c.Trello.AppKey
		credentials.Trello.AppKey = c.Trello.Token
	}

	return nil
}

func Reset() {
	credentials = nil
}

func getEnv(name string) string {
	return os.Getenv(name)
}
