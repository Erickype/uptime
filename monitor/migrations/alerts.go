package migrations

import "encore.app/site"

// SiteTransitionEvent defines the transition of a monitored site
// from up->down or down->up
type SiteTransitionEvent struct {
	// Site the site to be monitored
	Site *site.Site `json:"site,omitempty"`
	// Up specifies whether the site is up or down
	Up bool `json:"up,omitempty"`
}
