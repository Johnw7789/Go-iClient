package icloud

import (
	"fmt"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/google/uuid"
)

// * Generate a new frameId and clientId, then init the session
func (c *Client) AuthStart() (err error) {
	c.frameId = strings.ToLower(uuid.New().String())
	c.clientId = OAuthClientID

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoints[authStart], c.frameId, c.frameId, c.clientId, c.frameId), nil)
	if err != nil {
		return err
	}

	req.Header.Set(HdrAccept, "*/*")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	c.authAttr = resp.Header.Get("X-Apple-Auth-Attributes")

	return nil
}
