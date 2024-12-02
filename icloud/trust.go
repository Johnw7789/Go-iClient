package icloud

import (
	"fmt"

	http "github.com/bogdanfinn/fhttp"
)

// * Gets the trust and auth tokens, allowing for completing user authentication for iCloud web services
func (c *Client) GetTrust() (err error) {
	req, err := http.NewRequest(http.MethodGet, endpoints[trust], nil)
	if err != nil {
		return err
	}

	req.Header = c.updateRequestHeaders(req.Header.Clone())

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	c.authToken = resp.Header.Get("X-Apple-Session-Token")
	c.trustToken = resp.Header.Get("X-Apple-TwoSV-Trust-Token")

	return nil
}
