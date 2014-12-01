package github

import "time"

type Person struct {
	Username string
	Name     string
	Email    string
}

type PushEvent struct {
	Ref    string
	Pusher *Person

	Sender struct {
		Type     string
		ID       int
		Login    string
		Url      string
		HTML_URL string
	}

	Repository struct {
		ID       int
		Url      string
		HTML_URL string
		Git_URL  string
	}

	Head_Commit struct {
		ID        string
		Author    *Person
		Committer *Person
		Message   string
		Timestamp time.Time
		URL       string

		Added    []string // files added
		Modified []string // files modified
		Removed  []string // files removed
	}
}
