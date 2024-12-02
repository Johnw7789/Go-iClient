package icloud

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/tidwall/gjson"
)

// * Authenticate the user with the HME service using the trust and auth tokens
func (c *Client) AuthenticateHME() error {
	body := bytes.NewReader([]byte(fmt.Sprintf(`{"dsWebAuthToken":"%s","accountCountryCode":"USA","extended_login":true,"trustToken":"%s"}`, c.authToken, c.trustToken)))

	req, err := http.NewRequest(http.MethodPost, endpoints[hmeAuth], body)
	if err != nil {
		return err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// for each cookie under idmsa.apple.com, set it for p52-maildomainws.icloud.com as well
	u, _ := url.Parse("https://idmsa.apple.com")
	u2, _ := url.Parse("https://p52-maildomainws.icloud.com")

	cookies := c.HttpClient.GetCookies(u)
	c.HttpClient.SetCookies(u2, cookies)

	return nil
}

func (c *Client) generateHME() (string, error) {
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

func (c *Client) reserveHME(identifier, email string) (string, error) {
	randLabel := identifier + "-" + randomNumString()

	body := bytes.NewReader([]byte(`{"hme":"` + email + `","label":"` + randLabel + `","note":""}`))

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

// * Random num string between 1 and 8 length
func randomNumString() string {
	charSet := "1234567890"
	length := 8
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(charSet[rand.Intn(len(charSet))])
	}

	return result.String()
}

// * Generate a new HME and reserve it, if successful return the email to user
func (c *Client) GenerateHME(identifier string) (string, error) {
	hme, err := c.generateHME()
	if err != nil {
		return "", err
	}

	return c.reserveHME(identifier, hme)
}
