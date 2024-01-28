package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/partyhall/partyhall/config"
)

type spotifyAuth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`

	ExpiresAt time.Time `json:"-"`
}

func (sa *spotifyAuth) IsExpired() bool {
	return time.Now().After(sa.ExpiresAt.Add(-30 * time.Second))
}

type SpotifyImage struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type SpotifyTrack struct {
	Name    string `json:"name"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Album struct {
		Images []SpotifyImage `json:"images"`
	} `json:"album"`
}

func (s SpotifyTrack) String() string {
	return fmt.Sprintf("%v by %v", s.Name, s.Artists[0].Name)
}

type SpotifySearchResponse struct {
	Tracks struct {
		Items []SpotifyTrack `json:"items"`
	} `json:"tracks"`
}

type Spotify struct {
	auth *spotifyAuth
}

func (s *Spotify) authenticate() error {
	if s.auth != nil && !s.auth.IsExpired() {
		return nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	r, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(config.GET.SpotifyClientID+":"+config.GET.SpotifyClientSecret)))

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to authenticate: status code %v:\n%v", resp.StatusCode, err)
		}

		return fmt.Errorf("failed to authenticate: status code %v:\n%v", resp.StatusCode, string(body))
	}

	sa := spotifyAuth{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &sa)
	if err != nil {
		return err
	}

	sa.ExpiresAt = time.Now().Add(time.Duration(sa.ExpiresIn) * time.Second)
	s.auth = &sa

	return nil
}

func (s *Spotify) SearchSong(query string) ([]SpotifyTrack, error) {
	err := s.authenticate()
	if err != nil {
		return nil, err
	}

	r, _ := http.NewRequest("GET", "https://api.spotify.com/v1/search", nil)
	r.Header.Add("Authorization", "Bearer "+s.auth.AccessToken)

	q := r.URL.Query()
	q.Add("type", "track")
	q.Add("q", query)
	r.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to search: status code %v:\n%v", resp.StatusCode, err)
		}

		return nil, fmt.Errorf("failed to search: status code %v:\n%v", resp.StatusCode, string(body))
	}

	var searchResponse SpotifySearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	return searchResponse.Tracks.Items, err
}
