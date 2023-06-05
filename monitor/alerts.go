package monitor

import (
	"context"
	"encore.app/site"
	"encore.dev/pubsub"
	"encore.dev/storage/sqldb"
	"errors"
)

// SiteTransitionEvent defines the transition of a monitored site
// from up->down or down->up
type SiteTransitionEvent struct {
	// Site the site to be monitored
	Site *site.Site `json:"site,omitempty"`
	// Up specifies whether the site is up or down
	Up bool `json:"up,omitempty"`
}

var TransitionTopic = pubsub.NewTopic[*SiteTransitionEvent]("uptime-transition", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

// publishOnTransition gets the previous state and then decides to
// publish a TransitionTopic based whether the state has changed or not
func publishOnTransition(ctx context.Context, site *site.Site, isUp bool) error {
	wasUp, err := getPreviousSiteState(ctx, site.ID)
	if err != nil {
		return err
	}
	if wasUp == isUp {
		return nil
	}
	_, err = TransitionTopic.Publish(ctx, &SiteTransitionEvent{
		Site: site,
		Up:   isUp,
	})
	return err
}

// getPreviousSiteState returns the previous state of a site, if no previous ping then returns true
func getPreviousSiteState(ctx context.Context, siteID int) (up bool, err error) {
	query := "select up from checks where site_id = $1 order by checked_at limit 1"
	err = sqldb.QueryRow(ctx, query, siteID).Scan(&up)
	// check if the result has a row
	if errors.Is(err, sqldb.ErrNoRows) {
		// There was no previous ping; treat as if the site was up
		return true, nil
	}
	// error scanning the value
	if err != nil {
		return false, err
	}
	// return the previous state
	return up, nil
}
