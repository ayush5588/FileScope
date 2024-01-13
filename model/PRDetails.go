package model

// PR ...
type PR struct {
	Number       string   `json:"number"`
	URL          string   `json:"html_url"`
	State        string   `json:"state"`
	Title        string   `json:"title"`
	CreatedAt    string   `json:"created_at"`
	FilesChanged []string `json:"filesChanged"`
}
