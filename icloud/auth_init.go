package icloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	http "github.com/bogdanfinn/fhttp"
)

func (c *Client) AuthInit(A, email string) (authInitResp, error) {
	var body io.Reader

	reqBody := authInitReq{
		A:           A,
		AccountName: email,
		Protocols: []string{
			"s2k",
			"s2k_fo",
		},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return authInitResp{}, fmt.Errorf("marshal request body: %w", err)
	}

	body = bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, endpoints[authInit], body)
	if err != nil {
		return authInitResp{}, err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header = c.updateRequestHeaders(req.Header.Clone())

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return authInitResp{}, err
	}

	var authInitResp authInitResp
	if err := json.NewDecoder(resp.Body).Decode(&authInitResp); err != nil {
		return authInitResp, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return authInitResp, nil
}
