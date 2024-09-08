package discogs

import (
	"context"
	"fmt"
	"encoding/json"
)

type IdentityService service

type Identity struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	ResourceURL  string `json:"resource_url"`
	ConsumerName string `json:"consumer_name"`
}

func (s *IdentityService) Get() (id *Identity, err error) {
	req, err := s.client.NewRequest("GET", "oauth/identity", nil)
	if err != nil {
		err = fmt.Errorf("error building auth check request: %w", err)
		return
	}

	resp, err := s.client.Do(context.TODO(), req)
	if err != nil {
		err = fmt.Errorf("error sending auth check request: %w", err)
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("got non-200 from auth check endpoint: %d", resp.StatusCode)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&id)
	if err != nil {
		err = fmt.Errorf("error decoding identity response body: %w", err)
		return
	}

	return
}
