package icloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	http "github.com/bogdanfinn/fhttp"
)

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

// * uses the client channel to wait for and then submit the 2FA code
func (c *Client) submitTwoFactor() (err error) {
	// verification code has already be pushed to devices
	var body io.Reader
	var codeType string

	otp := <-c.OtpChannel

	reqBody := twoFactorCodeFromPhoneRequest{
		SecurityCode: &securityCode{
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
