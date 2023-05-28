package site

// Site is the struct to be monitored
type Site struct {
	ID  int    `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}
