package meta

type Credentials struct {
	GitHub *GitHubCredential
	Docker *DockerCredential
}

type GitHubCredential struct {
	Token string
}

type DockerCredential struct {
	Username string
	Password string
}
