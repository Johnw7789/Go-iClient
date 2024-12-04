package icloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/google/uuid"
)

// * Generate a new frameId and clientId, then init the session
func (c *Client) authStart() (err error) {
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

// * Submit the user's email to Apple
func (c *Client) authFederate(email string) error {
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

// * Initiate the authentication process, get the salt and B from the server
func (c *Client) authInit(A, email string) (AuthInitResp, error) {
	var body io.Reader

	reqBody := AuthInitReq{
		A:           A,
		AccountName: email,
		Protocols: []string{
			"s2k",
			"s2k_fo",
		},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return AuthInitResp{}, fmt.Errorf("marshal request body: %w", err)
	}

	body = bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, endpoints[authInit], body)
	if err != nil {
		return AuthInitResp{}, err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header = c.updateRequestHeaders(req.Header.Clone())

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return AuthInitResp{}, err
	}

	var authInitResp AuthInitResp
	if err := json.NewDecoder(resp.Body).Decode(&authInitResp); err != nil {
		return authInitResp, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return authInitResp, nil
}

// * Complete the authentication process and send the SRP data, as well as email
func (c *Client) authComplete(email, C, M1, M2 string) error {
	var body io.Reader

	reqBody := AuthCompleteReq{
		AccountName: email,
		RememberMe:  true,
		TrustTokens: []string{},
		M1:          M1,
		C:           C,
		M2:          M2,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	body = bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, endpoints[authComplete], body)
	if err != nil {
		return err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header = c.updateRequestHeaders(req.Header.Clone())

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch code := resp.StatusCode; code {
	case 200:
		return nil

	case 403:
		return ErrIncorrectUsernamePassword

	case 401:
		return ErrSeverErrorOrInvalidCreds

	case 409:
		// what we typically want to see. This is a 2FA or 2SA challenge
		return c.handleTwoFactor(resp)

	case 412:
		return ErrRequiredPrivacyAck

	case 502:
		return errors.New("apple server error")

	default:
		return ErrUnexpectedSigninResponse
	}
}

// * handleTwoFactor handles Two Factor authentication
func (c *Client) handleTwoFactor(signinResp *http.Response) error {
	// extract `X-Apple-Id-Session-Id` and `scnt` from response
	c.sessionID = signinResp.Header.Get(HdrXAppleIDSessionID)
	c.scnt = signinResp.Header.Get(HdrScnt)

	// // get authentication options
	req, err := http.NewRequest(http.MethodGet, endpoints[authOptions], nil)
	if err != nil {
		return err
	}
	// set required headers
	req.Header = c.updateRequestHeaders(req.Header.Clone())

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.submitTwoFactor()
}

// * Uses the client OTP channel to wait for and then submit the 2FA code
func (c *Client) submitTwoFactor() (err error) {
	otp := <-c.OtpChannel

	var body io.Reader
	var codeType string

	reqBody := TwoFactorCodeFromPhoneRequest{
		SecurityCode: &SecurityCode{
			Code: otp,
		},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	codeType = "trusteddevice"
	body = bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(endpoints[submitSecurityCode], codeType), body)
	if err != nil {
		return err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header = c.updateRequestHeaders(req.Header.Clone())

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// var submitResp submitSecurityCodeResponse
	// err = json.NewDecoder(resp.Body).Decode(&submitResp)
	// if err != nil {
	// 	return fmt.Errorf("could not unmarshal response body: %w", err)
	// }

	// for _, svcErr := range submitResp.ServiceErrors {
	// 	if strings.Contains(svcErr.Message, "Incorrect") {
	// 		return errors.New("invalid verification code")
	// 	}
	// }

	return nil
}

// * Authenticate the user for the various icloud web services using the trust and auth tokens
func (c *Client) authenticateWeb() error {
	body := bytes.NewReader([]byte(fmt.Sprintf(`{"dsWebAuthToken":"%s","accountCountryCode":"USA","extended_login":true,"trustToken":"%s"}`, c.authToken, c.trustToken)))

	req, err := http.NewRequest(http.MethodPost, endpoints[authWeb], body)
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

	// for each cookie under idmsa.apple.com, set it for icloud.com as well
	u, _ := url.Parse("https://idmsa.apple.com")
	u2, _ := url.Parse("https://icloud.com")

	cookies := c.HttpClient.GetCookies(u)
	c.HttpClient.SetCookies(u2, cookies)

	return nil
}
