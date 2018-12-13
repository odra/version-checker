package meta

type Credentials struct {
	GitHub *GitHubCredential
	Trello *TrelloCredential
	Docker *DockerCredential
}

type GitHubCredential struct {
	Token string
}

type TrelloCredential struct {
	AppKey string
	Token  string
}

type DockerCredential struct {
	Username string
	Password string
}
