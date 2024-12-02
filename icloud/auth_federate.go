package icloud

import (
	"bytes"
	"fmt"
	"io"

	http "github.com/bogdanfinn/fhttp"
)

func (c *Client) AuthFederate(email string) error {
	var body io.Reader

	data := `{"accountName":"` + email + `","rememberMe":true}`
	body = bytes.NewReader([]byte(data))

	req, err := http.NewRequest(http.MethodPost, endpoints[authFederate], body)
	if err != nil {
		return err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header = c.updateRequestHeaders(req.Header.Clone())

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
