package site

import "context"

// ListResponse is the response struct for List method that returns a []*Site
type ListResponse struct {
	// sites is the object of type []*Site
	Sites []*Site `json:"sites,omitempty"`
}

// List returns the list of monitored sites as a ListResponse structure
//
//encore:api public method=GET path=/site
func (s *Service) List(ctx context.Context) (*ListResponse, error) {
	var sites []*Site
	if err := s.db.Find(&sites).Error; err != nil {
		return nil, err
	}
	return &ListResponse{Sites: sites}, nil
}
