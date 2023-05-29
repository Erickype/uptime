package site

import "context"

// Delete deletes a monitored Site by passing its ID
//
//encore:api public method=DELETE path=/site/:siteID
func (s *Service) Delete(_ context.Context, siteID int) error {
	return s.db.Delete(&Site{ID: siteID}).Error
}
