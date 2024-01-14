package model

// CreationDate ...
type CreationDate struct {
	Day      string `json:"day"`
	Month    string `json:"month"`
	Year     string `json:"year"`
	FullDate string `json:"fullDate"`
}

// PR ...
type PR struct {
	Number    string       `json:"number"`
	URL       string       `json:"html_url"`
	State     string       `json:"state"`
	Title     string       `json:"title"`
	Branch    string       `json:"branch"`
	CreatedBy string       `json:"created_by"`
	CreatedOn CreationDate `json:"created_at"`
}
