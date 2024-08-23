package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

type CollectionService service

type Folder struct {
	ID          int    `json:"id"`
	Count       int    `json:"count"`
	Name        string `json:"name"`
	ResourceURL string `json:"resource_url"`
}

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

type releaseFormat struct {
	Name         string   `json:"name"`
	Quantity     int      `json:"qty,string"`
	Descriptions []string `json:"descriptions"`
}
type releaseBasicInfo struct {
	ID          int             `json:"id"`
	ResourceURL string          `json:"resource_url"`
	MasterID    int             `json:"master_id"`
	MasterURL   string          `json:"master_url"`
	Thumb       string          `json:"thumb"`
	CoverImage  string          `json:"cover_image"`
	Title       string          `json:"title"`
	Year        int             `json:"year"`
	Formats     []releaseFormat `json:"formats"`
	Artists     []Artist        `json:"artists"`
	Labels      []Label         `json:"labels"`
	Genres      []string        `json:"genres"`
	Styles      []string        `json:"styles"`
}

type Release struct {
	ID         int              `json:"id"`
	InstanceID int              `json:"instance_id"`
	DateAdded  string           `json:"date_added"`
	FolderID   int              `json:"folder_id"`
	Rating     int              `json:"rating"`
	BasicInfo  releaseBasicInfo `json:"basic_information"`
}

func (f *Folder) URL() (*url.URL, error) {
	return url.Parse(f.ResourceURL)
}

type GetFoldersResponse struct {
	Folders []Folder `json:"folders"`
}

func (s *CollectionService) GetFolders(ctx context.Context, username string) ([]Folder, error) {
	u := fmt.Sprintf("users/%s/collection/folders", username)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var r GetFoldersResponse
	_, err = s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	return r.Folders, nil
}

type GetReleaseByFolderResponse struct {
	// Paginator paginator `json:"pagination"`
	Releases []Release `json:"releases"`
}

func (s *CollectionService) GetReleasesByFolder(ctx context.Context, username string, folderID int) ([]Release, error) {
	var releases []Release

	u := fmt.Sprintf("users/%s/collection/folders/%d/releases", username, folderID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	first, pager, err := NewPager[GetReleaseByFolderResponse](resp, s.client)
	if err != nil {
		return nil, err
	}
	releases = append(releases, first.Releases...)

	var next *GetReleaseByFolderResponse
	for {
		next, err = pager.Next(ctx)
		if errors.Is(err, ErrPageDone) {
			break
		}
		if err != nil {
			return nil, err
		}
		releases = append(releases, next.Releases...)
	}

	return releases, nil
}
