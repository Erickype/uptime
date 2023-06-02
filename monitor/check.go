package monitor

import (
	"context"
	"encore.app/site"
	"encore.dev/storage/sqldb"
	"golang.org/x/sync/errgroup"
)

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

func check(ctx context.Context, site *site.Site) error {
	// Ping to the site
	response, err := Ping(ctx, site.URL)
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
