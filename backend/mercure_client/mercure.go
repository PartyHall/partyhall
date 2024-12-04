package mercure_client

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/state"
	"github.com/partyhall/partyhall/utils"
)

var CLIENT *Client

type Client struct {
	HubUrl     string
	Token      string
	HttpClient *http.Client
}

func NewClient() (*Client, error) {
	_, port, err := net.SplitHostPort(config.GET.ListeningAddr)
	if err != nil {

		return nil, err
	}

	hubURL := fmt.Sprintf("http://localhost:%v/.well-known/mercure", port)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.GetClaimsPublisher())
	tokenString, err := token.SignedString(config.GET.Mercure.PublisherKey)
	if err != nil {

		return nil, err
	}

	CLIENT = &Client{
		HubUrl:     hubURL,
		Token:      tokenString,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
	}

	return CLIENT, nil
}

func (c *Client) PublishEventWithType(topic string, data any, msgType string) error {
	if msgType == "" {
		msgType = topic
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Add("topic", topic)
	form.Add("type", msgType)
	form.Add("private", "true")
	form.Add("data", string(dataBytes))

	req, err := http.NewRequest("POST", c.HubUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	// Add required headers
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to publish event, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) PublishEvent(topic string, data any) error {
	return c.PublishEventWithType(topic, data, "")
}

func (c *Client) PublishSyncInProgress() {
	CLIENT.PublishEvent(
		"/sync-progress",
		map[string]any{
			"syncInProgress": state.STATE.SyncInProgress,
		},
	)
}
