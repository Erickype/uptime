package monitor

import (
	"context"
	"encore.app/site"
	"encore.dev/cron"
	"encore.dev/storage/sqldb"
	"golang.org/x/sync/errgroup"
)

// Check all tracked sites every n minutes
var _ = cron.NewJob("check-all", cron.JobConfig{
	Title:    "Check all sites",
	Endpoint: CheckAll,
	Every:    5 * cron.Minute,
})

// CheckAll checks all the availability of all sites
//
//encore:api public method=GET path=/checkAll
func CheckAll(ctx context.Context) error {
	listResponse, err := site.List(ctx)
	if err != nil {
		return err
	}

	// Check up to n sites concurrently based g.SetLimit
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(4)

	for _, webSite := range listResponse.Sites {
		Site := webSite
		g.Go(func() error {
			return check(ctx, Site)
		})
	}
	return g.Wait()
}

// Check checks the availability of a site based on its ID
//
//encore:api public method=GET path=/check/:siteID
func Check(ctx context.Context, siteID int) error {
	// Get the site
	Site, err := site.Get(ctx, siteID)
	if err != nil {
		return err
	}

	// Check the site
	return check(ctx, Site)
}

// check make a Ping request and then publish a transition based on the state
// finally inserts the request in checks database
func check(ctx context.Context, site *site.Site) error {
	// Ping to the site
	response, err := Ping(ctx, site.URL)
	if err != nil {
		return err
	}

	// Publish a pub/sub message if the site transitions from up->down or down->up
	err = publishOnTransition(ctx, site, response.Up)
	if err != nil {
		return err
	}

	// Insert the checked process
	query := "insert into public.checks (site_id, up, checked_at) values ($1,$2,now())"
	_, err = sqldb.Exec(ctx, query, site.ID, response.Up)
	if err != nil {
		return err
	}
	return nil
}
