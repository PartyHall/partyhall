package nexus

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (ns *NexusSync) fetchBackdropAlbumPage(page int) ([]ApiBackdropAlbum, int, error) {
	data, err := ns.doJsonRequest(http.MethodGet, fmt.Sprintf("/api/backdrop_albums?page=%v", page), nil)
	if err != nil {
		return nil, -1, err
	}

	albums := []ApiBackdropAlbum{}

	itemsRaw, ok := data["member"]
	if !ok {
		return nil, -1, errors.New("bad request from the api: no member in the response")
	}

	itemsSlice, ok := itemsRaw.([]interface{})
	if !ok {
		return nil, -1, errors.New("bad request from the api: member is not an array")
	}

	for _, item := range itemsSlice {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return nil, -1, fmt.Errorf("failed to marshal song data: %w", err)
		}

		var album ApiBackdropAlbum
		if err := json.Unmarshal(itemJSON, &album); err != nil {
			return nil, -1, fmt.Errorf("failed to unmarshal song: %w", err)
		}

		albums = append(albums, album)
	}

	totalItemsRaw, ok := data["totalItems"]
	if !ok {
		return nil, -1, errors.New("total items field is missing from response")
	}

	totalItemsFloat, ok := totalItemsRaw.(float64)
	if !ok {
		return nil, -1, errors.New("total items is not a number")
	}

	totalItems := int(totalItemsFloat)

	return albums, totalItems, nil
}

func (ns *NexusSync) fetchAllBackdropAlbums() ([]ApiBackdropAlbum, error) {
	total := -1

	albums := []ApiBackdropAlbum{}

	page := 1
	for total < 0 || len(albums) < total {
		pageAlbums, amtAlbums, err := ns.fetchBackdropAlbumPage(page)
		if err != nil {
			return nil, err
		}

		total = int(amtAlbums)
		albums = append(albums, pageAlbums...)
		page++
	}

	return albums, nil
}

func (ns *NexusSync) fetchBackdropPage(albumId int64, page int) ([]ApiBackdrop, int, error) {
	data, err := ns.doJsonRequest(http.MethodGet, fmt.Sprintf("/api/backdrop_albums/%v/backdrops?&page=%v", albumId, page), nil)
	if err != nil {
		return nil, -1, err
	}

	backdrops := []ApiBackdrop{}

	itemsRaw, ok := data["member"]
	if !ok {
		return nil, -1, errors.New("bad request from the api: no member in the response")
	}

	itemsSlice, ok := itemsRaw.([]interface{})
	if !ok {
		return nil, -1, errors.New("bad request from the api: member is not an array")
	}

	for _, item := range itemsSlice {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return nil, -1, fmt.Errorf("failed to marshal song data: %w", err)
		}

		var backdrop ApiBackdrop
		if err := json.Unmarshal(itemJSON, &backdrop); err != nil {
			return nil, -1, fmt.Errorf("failed to unmarshal song: %w", err)
		}

		backdrops = append(backdrops, backdrop)
	}

	totalItemsRaw, ok := data["totalItems"]
	if !ok {
		return nil, -1, errors.New("total items field is missing from response")
	}

	totalItemsFloat, ok := totalItemsRaw.(float64)
	if !ok {
		return nil, -1, errors.New("total items is not a number")
	}

	totalItems := int(totalItemsFloat)

	return backdrops, totalItems, nil
}

func (ns *NexusSync) fetchAllBackdrops(albumId int64) ([]ApiBackdrop, error) {
	total := -1

	backdrops := []ApiBackdrop{}

	page := 1
	for total < 0 || len(backdrops) < total {
		pageBackdrops, amtBackdrops, err := ns.fetchBackdropPage(albumId, page)
		if err != nil {
			return nil, err
		}

		total = int(amtBackdrops)
		backdrops = append(backdrops, pageBackdrops...)
		page++
	}

	return backdrops, nil
}
