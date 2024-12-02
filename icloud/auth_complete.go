package icloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	http "github.com/bogdanfinn/fhttp"
)

func (c *Client) AuthComplete(email, C, M1, M2 string) error {
	var body io.Reader

	reqBody := authCompleteReq{
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
		return errors.New("Apple server error")

	default:
		return ErrUnexpectedSigninResponse
	}
}
