package icloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	http "github.com/bogdanfinn/fhttp"
)

const (
	fmClientBuildNumber = "2602Build22"
	fmAppName           = "iCloud Find (Web)"
	fmAppVersion        = "2.0"
	fmAPIVersion        = "3.0"
)

type fmClientContext struct {
	AppName           string `json:"appName"`
	AppVersion        string `json:"appVersion"`
	APIVersion        string `json:"apiVersion"`
	DeviceListVersion int    `json:"deviceListVersion"`
	Fmly              bool   `json:"fmly"`
	Timezone          string `json:"timezone"`
	InactiveTime      int    `json:"inactiveTime"`
}

type fmRefreshClientReq struct {
	ServerContext json.RawMessage `json:"serverContext,omitempty"`
	ClientContext fmClientContext `json:"clientContext"`
}

type fmPlaySoundReq struct {
	ServerContext json.RawMessage `json:"serverContext"`
	ClientContext fmClientContext `json:"clientContext"`
	Device        string          `json:"device"`
	Subject       string          `json:"subject,omitempty"`
	UserAction    string          `json:"userAction"`
	Channels      []string        `json:"channels,omitempty"`
}

// GetDevices fetches the full list of Find My devices for the account (including family members).
// It must be called at least once before PlaySound to initialize the server context.
func (c *Client) GetDevices() ([]FMDevice, error) {
	resp, err := c.reqRefreshClient()
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

// GetDevice fetches all devices and returns the one matching the given ID.
// Returns an error if the device is not found.
func (c *Client) GetDevice(id string) (*FMDevice, error) {
	devices, err := c.GetDevices()
	if err != nil {
		return nil, err
	}
	for i := range devices {
		if devices[i].ID == id {
			return &devices[i], nil
		}
	}
	return nil, fmt.Errorf("device not found: %s", id)
}

// PlaySound triggers a sound alert on the device with the given ID.
//
// For AirPods and other multi-channel accessories, pass the audio channel names
// (e.g. []string{"left", "right"}). For all other devices (iPhone, Watch, Mac, etc.),
// pass nil or an empty slice.
//
// GetDevices must be called at least once before PlaySound to establish the server context.
// Returns the updated device state as reported by the server.
func (c *Client) PlaySound(deviceID string, channels []string) (*FMDevice, error) {
	return c.reqPlaySound(deviceID, channels)
}

// fmipURL builds the full URL for a fmipservice action including required query parameters.
func (c *Client) fmipURL(action string) string {
	return fmt.Sprintf(
		"%s/fmipservice/client/web/%s?clientBuildNumber=%s&clientMasteringNumber=%s&clientId=%s&dsid=%s",
		c.findMeURL, action, fmClientBuildNumber, fmClientBuildNumber, c.frameId, c.dsid,
	)
}

func (c *Client) defaultClientContext() fmClientContext {
	return fmClientContext{
		AppName:           fmAppName,
		AppVersion:        fmAppVersion,
		APIVersion:        fmAPIVersion,
		DeviceListVersion: 1,
		Fmly:              true,
		Timezone:          "US/Pacific",
		InactiveTime:      0,
	}
}

// serverContextWithID returns the stored server context with "id": "server_ctx" injected,
// ready to be echoed back to the fmipservice API.
func (c *Client) serverContextWithID() (json.RawMessage, error) {
	var ctxMap map[string]interface{}
	if err := json.Unmarshal(c.fmServerCtx, &ctxMap); err != nil {
		return nil, fmt.Errorf("unmarshal server context: %w", err)
	}
	ctxMap["id"] = "server_ctx"
	return json.Marshal(ctxMap)
}

func (c *Client) reqRefreshClient() (*FMDevicesResp, error) {
	if c.findMeURL == "" || c.dsid == "" {
		return nil, errors.New("find my not initialized: call Login() first")
	}

	reqBody := fmRefreshClientReq{
		ClientContext: c.defaultClientContext(),
	}

	// Echo back the previous server context if we have one.
	if c.fmServerCtx != nil && string(c.fmServerCtx) != "null" {
		ctx, err := c.serverContextWithID()
		if err != nil {
			return nil, err
		}
		reqBody.ServerContext = ctx
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.fmipURL("refreshClient"), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// fmipservice requires text/plain despite the body being JSON.
	req.Header.Set(HdrContentType, "text/plain;charset=UTF-8")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("referer", "https://www.icloud.com/")
	req.Header.Set("accept", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 450 {
		return nil, ErrFindMySessionExpired
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var devicesResp FMDevicesResp
	if err := json.Unmarshal(bodyBytes, &devicesResp); err != nil {
		return nil, err
	}

	// Cache the server context; it must be echoed back on all subsequent fmip calls.
	if devicesResp.ServerContext != nil && string(devicesResp.ServerContext) != "null" {
		c.fmServerCtx = devicesResp.ServerContext
	}

	return &devicesResp, nil
}

func (c *Client) reqPlaySound(deviceID string, channels []string) (*FMDevice, error) {
	if c.findMeURL == "" || c.dsid == "" {
		return nil, errors.New("find my not initialized: call Login() first")
	}
	if c.fmServerCtx == nil || string(c.fmServerCtx) == "null" {
		return nil, errors.New("server context unavailable: call GetDevices() before PlaySound()")
	}

	ctx, err := c.serverContextWithID()
	if err != nil {
		return nil, err
	}

	reqBody := fmPlaySoundReq{
		ServerContext: ctx,
		ClientContext: c.defaultClientContext(),
		Device:        deviceID,
		UserAction:    "PlaySound",
		Channels:      channels,
	}
	// The subject field is only sent for multi-channel accessories (AirPods).
	if len(channels) > 0 {
		reqBody.Subject = "Find My Accessory Alert"
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.fmipURL("playSound"), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set(HdrContentType, "text/plain;charset=UTF-8")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("referer", "https://www.icloud.com/")
	req.Header.Set("accept", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 450 {
		return nil, ErrFindMySessionExpired
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var soundResp FMDevicesResp
	if err := json.Unmarshal(bodyBytes, &soundResp); err != nil {
		return nil, err
	}

	if len(soundResp.Content) == 0 {
		return nil, errors.New("no device returned in play sound response")
	}

	return &soundResp.Content[0], nil
}
