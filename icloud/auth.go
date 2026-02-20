package icloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

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

// authComplete completes the authentication process and sends the SRP data.
func (c *Client) authComplete(email, C, M1, M2 string, otpProvider OTPProvider) error {
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
		return c.handleTwoFactor(resp, otpProvider)

	case 412:
		return ErrRequiredPrivacyAck

	case 502:
		return errors.New("apple server error")

	default:
		return ErrUnexpectedSigninResponse
	}
}

// handleTwoFactor handles Two Factor authentication.
func (c *Client) handleTwoFactor(signinResp *http.Response, otpProvider OTPProvider) error {
	// extract `X-Apple-Id-Session-Id` and `scnt` from response
	c.sessionID = signinResp.Header.Get(HdrXAppleIDSessionID)
	c.scnt = signinResp.Header.Get(HdrScnt)

	// get authentication options
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

	// Apple returns an updated scnt on every response; it must be forwarded on the next request.
	if newScnt := resp.Header.Get(HdrScnt); newScnt != "" {
		c.scnt = newScnt
	}

	return c.submitTwoFactor(otpProvider)
}

// submitTwoFactor calls the OTP provider and submits the 2FA code.
func (c *Client) submitTwoFactor(otpProvider OTPProvider) error {
	otp, err := otpProvider()
	if err != nil {
		return fmt.Errorf("failed to get OTP: %w", err)
	}

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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse dsid and service URLs for use in subsequent service calls.
	var accountResp struct {
		DsInfo struct {
			Dsid string `json:"dsid"`
		} `json:"dsInfo"`
		Webservices struct {
			Account struct {
				URL string `json:"url"`
			} `json:"account"`
			Findme struct {
				URL string `json:"url"`
			} `json:"findme"`
			Contacts struct {
				URL string `json:"url"`
			} `json:"contacts"`
		} `json:"webservices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&accountResp); err == nil {
		c.dsid = accountResp.DsInfo.Dsid
		c.accountURL = accountResp.Webservices.Account.URL
		c.findMeURL = accountResp.Webservices.Findme.URL
		c.contactsURL = accountResp.Webservices.Contacts.URL
	}

	// for each cookie under idmsa.apple.com, set it for icloud.com as well
	u, _ := url.Parse("https://idmsa.apple.com")
	u2, _ := url.Parse("https://icloud.com")

	cookies := c.HttpClient.GetCookies(u)
	c.HttpClient.SetCookies(u2, cookies)

	// Make a second accountLogin to the partition-specific setup server with the dsid.
	// This binds the iCloud session to the user's assigned partition and elevates
	// access for partition-local services like Find My. Without this step, fmipservice
	// returns 450 (re-auth required) even with a valid generic iCloud session.
	if c.accountURL != "" && c.dsid != "" {
		if err := c.authenticateWebPartition(); err != nil {
			// Non-fatal: basic iCloud session still works; Find My may require re-auth.
			_ = err
		}
	}

	return nil
}

// authenticateWebPartition performs a second accountLogin to the user's partition-specific
// setup server (e.g. p52-setup.icloud.com) with the dsid. This is required to elevate
// the session for partition-local services such as Find My (fmipservice).
func (c *Client) authenticateWebPartition() error {
	partitionURL := fmt.Sprintf(
		"%s/setup/ws/1/accountLogin?clientBuildNumber=2602Build17&clientMasteringNumber=2602Build17&clientId=%s&dsid=%s",
		c.accountURL, c.frameId, c.dsid,
	)

	body := bytes.NewReader([]byte(fmt.Sprintf(`{"dsWebAuthToken":"%s","accountCountryCode":"USA"}`, c.authToken)))

	req, err := http.NewRequest(http.MethodPost, partitionURL, body)
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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("partition accountLogin: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// ValidateSession checks whether the current iCloud session is still valid.
// Returns true if the session is active, false if it has expired.
// Call Login() again if this returns false.
func (c *Client) ValidateSession() (bool, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(endpoints[authValidate], c.frameId, c.dsid), nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("referer", "https://www.icloud.com/")
	req.Header.Set(HdrAccept, "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Success, nil
}

// KeepAlive keeps the iCloud session alive by calling ValidateSession on the
// given interval. It blocks until the context is cancelled or the session expires.
// Intended to be run in a goroutine:
//
//	go func() {
//	    if err := client.KeepAlive(ctx, 5*time.Minute); err != nil {
//	        // session expired — call Login() again
//	    }
//	}()
func (c *Client) KeepAlive(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			valid, err := c.ValidateSession()
			if err != nil {
				return err
			}
			if !valid {
				return ErrSessionExpired
			}
		}
	}
}
