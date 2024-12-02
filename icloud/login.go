package icloud

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"

	b64 "encoding/base64"

	"github.com/Johnw7789/go-icloud/srp"
)

// * LoginInit() handles the login process up to the point of trusting the device
func (c *Client) LoginInit() error {
	err := c.AuthStart()
	if err != nil {
		return err
	}

	err = c.AuthFederate(c.Username)
	if err != nil {
		return err
	}

	params := srp.GetParams(2048)
	params.NoUserNameInX = true // this is required for Apple's implementation

	client := srp.NewSRPClient(params, nil)

	// * Get the salt and B from the server
	authInitResp, err := c.AuthInit(b64.StdEncoding.EncodeToString(client.GetABytes()), c.Username)
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

	return c.AuthComplete(c.Username, authInitResp.C, b64.StdEncoding.EncodeToString(client.M1), b64.StdEncoding.EncodeToString(client.M2))
}

// * Login() logs in to iCloud and finishes authentication required for the HME service
func (c *Client) Login() error {
	err := c.LoginInit()
	if err != nil {
		return err
	}

	// * Trust the device
	err = c.GetTrust()
	if err != nil {
		return err
	}

	// * Auth required for the HME service
	return c.AuthenticateHME()
}
