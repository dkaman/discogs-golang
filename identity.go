package discogs

import (
	"context"
	"encoding/json"
	"fmt"
)

type IdentityService service

type Identity struct {
	ID           int    `json:"id"`
	ResourceURL  string `json:"resource_url"`
	Username     string `json:"username"`
	ConsumerName string `json:"consumer_name"`
}

type Profile struct {
	ID                   int     `json:"id"`
	ResourceURL          string  `json:"resource_url"`
	Profile              string  `json:"profile"`
	WantlistURL          string  `json:"wantlist_url"`
	Rank                 int     `json:"rank"`
	NumPending           int     `json:"num_pending"`
	NumForSale           int     `json:"num_for_sale"`
	HomePage             string  `json:"home_page"`
	Location             string  `json:"location"`
	CollectionFoldersURL string  `json:"collection_folders_url"`
	Username             string  `json:"username"`
	CollectionFieldsURL  string  `json:"collection_fields_url"`
	ReleasesContributed  int     `json:"releases_contributed"`
	Registered           string  `json:"registered"`
	RatingAvg            float64 `json:"rating_avg"`
	NumCollection        int     `json:"num_collection"`
	ReleasesRated        int     `json:"releases_rated"`
	NumLists             int     `json:"num_lists"`
	Name                 string  `json:"name"`
	NumWantlist          int     `json:"num_wantlist"`
	InventoryURL         string  `json:"inventory_url"`
	AvatarURL            string  `json:"avatar_url"`
	BannerURL            string  `json:"banner_url"`
	URI                  string  `json:"uri"`
	BuyerRating          float64 `json:"buyer_rating"`
	BuyerRatingStars     int     `json:"buyer_rating_stars"`
	BuyerNumRatings      int     `json:"buyer_num_ratings"`
	SellerRating         float64 `json:"seller_rating"`
	SellerRatingStars    int     `json:"seller_rating_stars"`
	SellerNumRatings     int     `json:"seller_num_ratings"`
	CurrAbbr             string  `json:"curr_abbr"`
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

func (s *IdentityService) GetProfile(username string) (profile *Profile, err error) {
	return
}

func (s *IdentityService) EditProfile(username string) (profile *Profile, err error) {
	return
}

func (s *IdentityService) GetSubmissions(username string) (err error) {
	return
}

func (s *IdentityService) GetContributions(username string, sort string, sortOrder string) (err error) {
	return
}
