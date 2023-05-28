package site

import "context"

// Site is the struct to be monitored
type Site struct {
	ID  int    `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}

// Get gets a Site by its ID.
//
//encore:api public method=GET path=/site/:siteID
func (s *Service) Get(_ context.Context, siteID int) (*Site, error) {
	var site Site
	if err := s.db.Where("id = $1", siteID).First(&site).Error; err != nil {
		return nil, err
	}
	return &site, nil
}
