package monitor

import (
	"context"
	"encore.app/site"
	"encore.dev/storage/sqldb"
)

// Check checks the availability of a site based on its ID
//
//encore:api public method=GET path=/check/:siteID
func Check(ctx context.Context, siteID int) error {
	// Get the site
	Site, err := site.Get(ctx, siteID)
	if err != nil {
		return err
	}
	// Ping to the site
	response, err := Ping(ctx, Site.URL)
	if err != nil {
		return err
	}
	// Insert the checked process
	query := "insert into public.checks (site_id, up, checked_at) values ($1,$2,now())"
	_, err = sqldb.Exec(ctx, query, Site.ID, response.Up)
	if err != nil {
		return err
	}
	return nil
}
