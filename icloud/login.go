package icloud

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/pbkdf2"

	b64 "encoding/base64"

	"github.com/Johnw7789/Go-iClient/srp"
	http "github.com/bogdanfinn/fhttp"
)

// * Login() logs in to iCloud and authenticates the user to access any iCloud web app/service
func (c *Client) Login() error {
	err := c.loginInit()
	if err != nil {
		return err
	}

	err = c.getTrust()
	if err != nil {
		return err
	}

	return c.authenticateWeb()
}

// * loginInit() handles the login process up to the point of trusting the device
func (c *Client) loginInit() error {
	err := c.authStart()
	if err != nil {
		return err
	}

	err = c.authFederate(c.Username)
	if err != nil {
		return err
	}

	params := srp.GetParams(2048)
	params.NoUserNameInX = true // this is required for Apple's implementation

	client := srp.NewSRPClient(params, nil)

	// * Get the salt and B from the server
	authInitResp, err := c.authInit(b64.StdEncoding.EncodeToString(client.GetABytes()), c.Username)
	if err != nil {
		return err
	}

	// * Both the salt and B are base64 encoded so we need to decode them
	bDec, err := b64.StdEncoding.DecodeString(authInitResp.B)
	if err != nil {
		return err
	}

	saltDec, err := b64.StdEncoding.DecodeString(authInitResp.Salt)
	if err != nil {
		return err
	}

	// * Generate the password key
	passHash := sha256.Sum256([]byte(c.Password))
	passKey := pbkdf2.Key(passHash[:], []byte(saltDec), authInitResp.Iteration, 32, sha256.New)

	// * Process the challenge using the server provided salt and B
	client.ProcessClientChanllenge([]byte(c.Username), passKey, saltDec, bDec)

	return c.authComplete(c.Username, authInitResp.C, b64.StdEncoding.EncodeToString(client.M1), b64.StdEncoding.EncodeToString(client.M2))
}

// * getTrust() Gets the trust and auth tokens, allowing for completing user authentication for iCloud web services
func (c *Client) getTrust() (err error) {
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
