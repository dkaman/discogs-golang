package discogs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type DatabaseService service

const (
	Unknown = ""
	Correct = "Correct"
)

type Artist struct {
	ID          int    `json:"id"`
	ResourceURL string `json:"resource_url"`
	Name        string `json:"name"`
	Anv         string `json:"anv"`
	Join        string `json:"join"`
	Role        string `json:"role"`
	Tracks      string `json:"tracks"`
}

type Label struct {
	ID             int    `json:"id"`
	ResourceURL    string `json:"resource_url"`
	Name           string `json:"name"`
	CatNo          string `json:"catno"`
	EntityType     int    `json:"entity_type,string"`
	EntityTypeName string `json:"entity_type_name"`
}

type Company struct {
	ID             int    `json:"id"`
	ResourceURL    string `json:"resource_url"`
	Name           string `json:"name"`
	CatNo          string `json:"catno"`
	EntityType     int    `json:"entity_type,string"`
	EntityTypeName string `json:"entity_type_name"`
}

type User struct {
	ResourceURL string `json:"resource_url"`
	Username    string `json:"username"`
}

type Track struct {
	Duration string `json:"duration"`
	Position string `json:"position"`
	Title    string `json:"title"`
	Type     string `json:"type_"`
}

type Format struct {
	Descriptions []string `json:"descriptions"`
	Name         string   `json:"name"`
	Qty          string   `json:"qty"`
}

type Release struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Artists     []Artist `json:"artists"`
	DataQuality string   `json:"data_quality"`
	Thumb       string   `json:"thumb"`

	Community struct {
		Contributors []User `json:"contributors"`
	} `json:"community"`

	Companies       []Company `json:"companies"`
	Country         string    `json:"country"`
	DateAdded       string    `json:"date_added"`
	DateChanged     string    `json:"date_changed"`
	EstimatedWeight int       `json:"estimated_weight"`
	ExtraArtists    []Artist  `json:"extraartists"`
	FormatQuantity  int       `json:"format_quantity"`
	Formats         []Format  `json:"formats"`
	Genres          []string  `json:"genres"`

	Identifiers []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"intentifiers"`

	Images []struct {
		Height      int    `json:"height"`
		ResourceURL string `json:"resource_url"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
		URI150      string `json:"uri150"`
		Width       int    `json:"width"`
	} `json:"images"`

	Labels            []Label `json:"labels"`
	LowestPrice       float64 `json:"lowest_price"`
	MasterID          int     `json:"master_id"`
	MasterURL         string  `json:"master_url"`
	Notes             string  `json:"notes"`
	NumForSale        int     `json:"num_for_sale"`
	Released          string  `json:"released"`
	ReleasedFormatted string  `json:"released_formatted"`
	ResourceURL       string  `json:"resource_url"`

	// TODO what is this?
	Series []any `json:"series"`

	Status    string   `json:"status"`
	Styles    []string `json:"styles"`
	Tracklist []Track  `json:"tracklist"`

	URI string `json:"uri"`

	Videos []struct {
		Description string `json:"description"`
		Duration    int    `json:"duration"`
		Embed       bool   `json:"embed"`
		Title       string `json:"title"`
		URI         string `json:"uri"`
	} `json:"videos"`

	Year int `json:"year"`
}

type GetReleaseResonse struct {
	Release *Release
}

func (s *DatabaseService) GetRelease(ctx context.Context, id int) (release *Release, err error) {
	u := fmt.Sprintf("/releases/%d", id)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 from server, got %d", resp.StatusCode)
	}

	var out Release

	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return
	}

	release = &out

	return
}
