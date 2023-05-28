package site

import "context"

// AddParams is the struct to add a Site
type AddParams struct {
	// URL is the url of the Site.If it doesn't contain a scheme
	// (like "http:" or "https:") it defaults to "https:"."
	URL string `json:"url,omitempty"`
}

// Add adds a new Site to the list of monitored sites
//
//encore:api public method=POST path=/site
func (s *Service) Add(_ context.Context, params *AddParams) (*Site, error) {
	site := &Site{URL: params.URL}
	if err := s.db.Create(site).Error; err != nil {
		return nil, err
	}
	return site, nil
}
