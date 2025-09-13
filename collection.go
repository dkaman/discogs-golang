package discogs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type CollectionService service

type Folder struct {
	ID          int    `json:"id"`
	ResourceURL string `json:"resource_url"`
	Count       int    `json:"count"`
	Name        string `json:"name"`
}

type ReleaseInstance struct {
	ID         int    `json:"id"`
	InstanceID int    `json:"instance_id"`
	DateAdded  string `json:"date_added"`
	FolderID   int    `json:"folder_id"`
	Rating     int    `json:"rating"`
	BasicInfo  struct {
		ID          int      `json:"id"`
		ResourceURL string   `json:"resource_url"`
		MasterID    int      `json:"master_id"`
		MasterURL   string   `json:"master_url"`
		Thumb       string   `json:"thumb"`
		CoverImage  string   `json:"cover_image"`
		Title       string   `json:"title"`
		Year        int      `json:"year"`
		Genres      []string `json:"genres"`
		Styles      []string `json:"styles"`
		Artists     []Artist `json:"artists"`
		Labels      []Label  `json:"labels"`
		Formats     []Format `json:"formats"`
	} `json:"basic_information"`
}

type Field struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Options  []string `json:"options"`
	Position int      `json:"position"`
	Type     string   `json:"type"`
	Public   bool     `json:"public"`
	Lines    int      `json:"lines"`
}

type Value struct {
	Maximum string `json:"maximum"`
	Median  string `json:"median"`
	Minimum string `json:"minimum"`
}

type Instance struct {
	ID          int    `json:"instance_id"`
	ResourceURL string `json:"resource_url"`
}

type releaseResponse struct {
	Releases []ReleaseInstance `json:"releases"`
}

func (f *Folder) URL() (*url.URL, error) {
	return url.Parse(f.ResourceURL)
}

type GetFoldersResponse struct {
	Folders []Folder `json:"folders"`
}

func (s *CollectionService) ListFolders(ctx context.Context, username string) (folders []Folder, err error) {
	u := fmt.Sprintf("users/%s/collection/folders", username)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return
	}

	var r GetFoldersResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return
	}

	folders = r.Folders

	return
}

func (s *CollectionService) CreateFolder(ctx context.Context, username string, folderName string) (folder *Folder, err error) {
	u := fmt.Sprintf("users/%s/collection/folders", username)

	body := struct {
		Username string `json:"username"`
		Name     string `json:"name"`
	}{
		Username: username,
		Name:     folderName,
	}

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("did not recieve 201 from server, got %d", resp.StatusCode)
	}

	var f Folder
	err = json.NewDecoder(resp.Body).Decode(&f)
	folder = &f

	return
}

func (s *CollectionService) GetFolder(ctx context.Context, username string, folderID int) (folder *Folder, err error) {
	u := fmt.Sprintf("users/%s/collection/folders/%d", username, folderID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("did not receive 200 from server, got %d", resp.StatusCode)
	}

	var f Folder
	err = json.NewDecoder(resp.Body).Decode(&f)
	folder = &f

	return
}

func (s *CollectionService) EditFolder(ctx context.Context, username string, folderID int, newFolder Folder) (folder *Folder, err error) {
	u := fmt.Sprintf("users/%s/collection/folders/%d", username, folderID)

	req, err := s.client.NewRequest(http.MethodPost, u, newFolder)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 from server, got %d", resp.StatusCode)
	}

	var f Folder
	err = json.NewDecoder(resp.Body).Decode(&f)
	folder = &f

	return
}

func (s *CollectionService) DeleteFolder(ctx context.Context, username string, folderID int) (err error) {
	u := fmt.Sprintf("users/%s/collection/folders/%d", username, folderID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent { // 204
		err = fmt.Errorf("obtained non-204 on delete: got %d", resp.StatusCode)
	}

	return
}

type GetFolderByReleaseResponse releaseResponse

func (s *CollectionService) GetFolderByRelease(ctx context.Context, username string, releaseID int) (releases []ReleaseInstance, err error) {
	u := fmt.Sprintf("users/%s/collection/releases/%d", username, releaseID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 from server, got %d", resp.StatusCode)
	}

	first, pager, err := NewPager[GetFolderByReleaseResponse](resp, s.client)
	if err != nil {
		return nil, err
	}
	releases = append(releases, first.Releases...)

	for {
		next, err := pager.Next(ctx)
		if errors.Is(err, ErrPageDone) {
			break
		}
		if err != nil {
			return nil, err
		}
		releases = append(releases, next.Releases...)
	}

	return
}

type GetReleaseByFolderResponse releaseResponse

func (s *CollectionService) GetReleasesByFolder(ctx context.Context, username string, folderID int) (releases []ReleaseInstance, err error) {
	u := fmt.Sprintf("users/%s/collection/folders/%d/releases", username, folderID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 from server, got %d", resp.StatusCode)
	}

	first, pager, err := NewPager[GetReleaseByFolderResponse](resp, s.client)
	if err != nil {
		return nil, err
	}
	releases = append(releases, first.Releases...)

	for {
		next, err := pager.Next(ctx)
		if errors.Is(err, ErrPageDone) {
			break
		}
		if err != nil {
			return nil, err
		}
		releases = append(releases, next.Releases...)
	}

	return
}

type AddReleaseToFolderResponse Instance

func (s *CollectionService) AddReleaseToFolder(ctx context.Context, username string, folderID int, releaseID int) (instance Instance, err error) {
	u := fmt.Sprintf("users/%s/collection/folders/%d/releases/%d", username, folderID, releaseID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusCreated { // 201
		err = fmt.Errorf("did not receive status 201 from server, got %d", resp.StatusCode)
		return
	}

	var r AddReleaseToFolderResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return
	}

	instance = Instance(r)

	return
}

func (s *CollectionService) ChangeRatingOfRelease(ctx context.Context, username string, folderID int, releaseID int, instanceID int, rating int) (err error) {
	u := fmt.Sprintf("users/%s/collection/folders/%d/releases/%d/instances/%d", username, folderID, releaseID, instanceID)

	body := struct {
		Rating int `json:"rating"`
	}{
		Rating: rating,
	}

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent { // 204
		err = fmt.Errorf("did not receive status 204 from server, got %d", resp.StatusCode)
		return
	}

	return
}

func (s *CollectionService) RemoveReleaseFromFolder(ctx context.Context, username string, folderID int, releaseID int, instanceID int) (err error) {
	u := fmt.Sprintf("users/%s/collection/folders/%d/releases/%d/instances/%d", username, folderID, releaseID, instanceID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return
	}

	switch resp.StatusCode {
	case http.StatusNoContent: // 204
		err = nil
	case http.StatusForbidden: // 403
		err = fmt.Errorf("unauthorized to edit folder %d", folderID)
	default:
		err = fmt.Errorf("did not receive 204 from server, got %d", resp.StatusCode)
	}

	return
}

type ListCustomFieldsResponse struct {
	Fields []Field `json:"fields"`
}

func (s *CollectionService) ListCustomFields(ctx context.Context, username string) (fields []Field, err error) {
	u := fmt.Sprintf("users/%s/collection/fields", username)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 from server, got %d", resp.StatusCode)
	}

	var r ListCustomFieldsResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	fields = r.Fields

	return
}

func (s *CollectionService) EditCustomFields(ctx context.Context, username string, folderID int, releaseID int, instanceID int, fieldID int, value string) (err error) {
	u := fmt.Sprintf("users/%s/collection/folders/%d/releases/%d/instances/%d/fields/%d", username, folderID, releaseID, instanceID, fieldID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent { // 204
		err = fmt.Errorf("did not receive 204 from server, got %d", resp.StatusCode)
	}

	return
}

func (s *CollectionService) GetCollectionValue(ctx context.Context, username string) (value *Value, err error) {
	u := fmt.Sprintf("users/%s/collection/value", username)

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

	var r Value
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	value = &r

	return
}
