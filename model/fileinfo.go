package model

// FileInfo represents components of path to file
type FileInfo struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	Path   string `json:"path"`
	URL    string `json:"url"`
}
