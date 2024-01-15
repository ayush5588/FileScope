package model

type Core struct {
	Limit     int `json:"limit"`
	Used      int `json:"used"`
	Remaining int `json:"remaining"`
	Reset     int `json:"reset"`
}

type Resources struct {
	Core Core
}

type Rate struct {
	Limit     int `json:"limit"`
	Used      int `json:"used"`
	Remaining int `json:"remaining"`
	Reset     int `json:"reset"`
}

// RateLimitResBody ...
type RateLimitResBody struct {
	Resources Resources `json:"resources"`
	Rate      Rate      `json:"rate"`
}
