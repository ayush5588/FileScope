package model

// FileInfo represents components of path to file
type FileInfo struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	Path   string `json:"path"`
	URL    string `json:"url"`
}

// File represents details of the file in a PR
type File struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	SHA    string `json:"sha"`
}
