package icloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/tidwall/gjson"
)

// * RetrieveHMEList() Retrieves a list of HMEs from the user's account
func (c *Client) RetrieveHMEList() ([]HmeEmail, error) {
	resp, err := c.reqRetrieveHMEList()
	if err != nil {
		return nil, err
	}

	var hmeListResp HMEListResp
	err = json.Unmarshal([]byte(resp), &hmeListResp)
	if err != nil {
		return nil, err
	}

	return hmeListResp.Result.HmeEmails, nil
}

// * ReserveHME() Generates a new HME and reserves it, if successful it returns the email to user
func (c *Client) ReserveHME(label, note string) (string, error) {
	hme, err := c.reqGenerateHME()
	if err != nil {
		return "", err
	}

	return c.reqReserveHME(hme, label, note)
}

// * DeactivateHME() Deactivates an HME from the user's account, using the given anonymousId. Can be reactivated later
func (c *Client) DeactivateHME(anonymousId string) (bool, error) {
	resp, err := c.reqDeactivateHME(anonymousId)
	if err != nil {
		return false, err
	}

	return gjson.Get(resp, "success").Bool(), nil
}

// * ReactivateHME() Reactivates an HME for a user's account, using the given anonymousId
func (c *Client) ReactivateHME(anonymousId string) (bool, error) {
	resp, err := c.reqReactivateHME(anonymousId)
	if err != nil {
		return false, err
	}

	return gjson.Get(resp, "success").Bool(), nil
}

// * DeleteHME() Deletes an HME from the user's account, using the given anonymousId. Can not be recovered later
func (c *Client) DeleteHME(anonymousId string) (bool, error) {
	resp, err := c.reqDeleteHME(anonymousId)
	if err != nil {
		return false, err
	}

	return gjson.Get(resp, "success").Bool(), nil
}

func (c *Client) reqRetrieveHMEList() (string, error) {
	req, err := http.NewRequest(http.MethodGet, endpoints[hmeList], nil)
	if err != nil {
		return "", err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (c *Client) reqGenerateHME() (string, error) {
	body := bytes.NewReader([]byte(`{"langCode":"en-us"}`))

	req, err := http.NewRequest(http.MethodPost, endpoints[hmeGen], body)
	if err != nil {
		return "", err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	success := gjson.Get(string(bodyBytes), "success").Bool()
	if !success {
		return "", errors.New("failed to generate HME")
	}

	email := gjson.Get(string(bodyBytes), "result.hme").String()

	return email, nil
}

func (c *Client) reqReserveHME(email, label, note string) (string, error) {
	body := bytes.NewReader([]byte(fmt.Sprintf(`{"hme":"%s","label":"%s","note":"%s"}`, email, label, note)))

	req, err := http.NewRequest(http.MethodPost, endpoints[hmeReserve], body)
	if err != nil {
		return "", err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !strings.Contains(string(bodyBytes), email) {
		return "", errors.New("failed to reserve HME")
	}

	return email, nil
}

func (c *Client) reqDeactivateHME(anonymousId string) (string, error) {
	body := bytes.NewReader([]byte(`{"anonymousId":"` + anonymousId + `"}`))

	req, err := http.NewRequest(http.MethodPost, endpoints[hmeDeactivate], body)
	if err != nil {
		return "", err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (c *Client) reqReactivateHME(anonymousId string) (string, error) {
	body := bytes.NewReader([]byte(`{"anonymousId":"` + anonymousId + `"}`))

	req, err := http.NewRequest(http.MethodPost, endpoints[hmeReactivate], body)
	if err != nil {
		return "", err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (c *Client) reqDeleteHME(anonymousId string) (string, error) {
	body := bytes.NewReader([]byte(`{"anonymousId":"` + anonymousId + `"}`))

	req, err := http.NewRequest(http.MethodPost, endpoints[hmeDelete], body)
	if err != nil {
		return "", err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
