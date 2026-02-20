package icloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/google/uuid"
)

const (
	coClientBuildNumber = "2604Build17"
	coClientVersion     = "2.1"
	coLocale            = "en_US"
	coOrder             = "last,first"
)

func (c *Client) buildContactsURL(action string) string {
	return fmt.Sprintf("%s/co/%s?clientBuildNumber=%s&clientId=%s&clientMasteringNumber=%s&clientVersion=%s&dsid=%s&locale=%s&order=%s", c.contactsURL, action, coClientBuildNumber, c.frameId, coClientBuildNumber, coClientVersion, c.dsid, coLocale, url.QueryEscape(coOrder))
}

func (c *Client) buildContactsURLWithTokens(action string) string {
	return fmt.Sprintf("%s&prefToken=%s&syncToken=%s", c.buildContactsURL(action), url.QueryEscape(c.prefToken), url.QueryEscape(c.syncToken))
}

// GetContacts fetches all contacts from iCloud. On first call, this automatically
// initializes the contacts session by calling startup. Returns the full contact list.
func (c *Client) GetContacts() ([]Contact, error) {
	if c.syncToken == "" || c.prefToken == "" {
		startupResp, err := c.reqStartup()
		if err != nil {
			return nil, err
		}
		return startupResp.Contacts, nil
	}

	if c.contactsURL == "" || c.dsid == "" {
		return nil, fmt.Errorf("contacts not available: missing contactsURL or dsid")
	}

	req, err := http.NewRequest(http.MethodGet, c.buildContactsURLWithTokens("contacts/card/"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("referer", "https://www.icloud.com/")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 450 {
		return nil, ErrSessionExpired
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get contacts failed with status %d: %s", resp.StatusCode, string(body))
	}

	var contactsResp ContactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&contactsResp); err != nil {
		return nil, fmt.Errorf("decode contacts response: %w", err)
	}

	if contactsResp.SyncToken != "" {
		c.syncToken = contactsResp.SyncToken
	}
	if contactsResp.PrefToken != "" {
		c.prefToken = contactsResp.PrefToken
	}

	return contactsResp.Contacts, nil
}

func (c *Client) GetContact(contactID string) (*Contact, error) {
	contacts, err := c.GetContacts()
	if err != nil {
		return nil, err
	}

	for _, contact := range contacts {
		if contact.ContactID == contactID {
			return &contact, nil
		}
	}

	return nil, fmt.Errorf("contact not found: %s", contactID)
}

// CreateContact creates a new contact in iCloud. A contactId UUID is generated
// automatically by the client. Returns the created contact with etag populated.
func (c *Client) CreateContact(contact Contact) (*Contact, error) {
	if err := c.ensureContactsInit(); err != nil {
		return nil, err
	}

	if c.contactsURL == "" || c.dsid == "" {
		return nil, fmt.Errorf("contacts not available: missing contactsURL or dsid")
	}

	// The API requires a client-generated contactId UUID
	if contact.ContactID == "" {
		contact.ContactID = strings.ToUpper(uuid.New().String())
	}

	reqBody := struct {
		Contacts []Contact `json:"contacts"`
	}{
		Contacts: []Contact{contact},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.buildContactsURLWithTokens("contacts/card/"), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("referer", "https://www.icloud.com/")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 450 {
		return nil, ErrSessionExpired
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create contact failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createResp ContactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("decode create response: %w", err)
	}

	if createResp.SyncToken != "" {
		c.syncToken = createResp.SyncToken
	}
	if createResp.PrefToken != "" {
		c.prefToken = createResp.PrefToken
	}

	if len(createResp.Contacts) == 0 {
		return nil, fmt.Errorf("create contact: no contact returned in response")
	}

	return &createResp.Contacts[0], nil
}

// UpdateContact updates an existing contact. Requires ContactID and Etag to be set.
// Returns the updated contact with new Etag.
func (c *Client) UpdateContact(contact Contact) (*Contact, error) {
	if contact.ContactID == "" {
		return nil, fmt.Errorf("update contact: ContactID is required")
	}
	if contact.Etag == "" {
		return nil, fmt.Errorf("update contact: Etag is required")
	}

	if err := c.ensureContactsInit(); err != nil {
		return nil, err
	}

	if c.contactsURL == "" || c.dsid == "" {
		return nil, fmt.Errorf("contacts not available: missing contactsURL or dsid")
	}

	reqBody := struct {
		Contacts []Contact `json:"contacts"`
	}{
		Contacts: []Contact{contact},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.buildContactsURLWithTokens("contacts/card/")+"&method=PUT", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("referer", "https://www.icloud.com/")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 450 {
		return nil, ErrSessionExpired
	}

	if resp.StatusCode == 409 {
		return nil, ErrContactEtagMismatch
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update contact failed with status %d: %s", resp.StatusCode, string(body))
	}

	var updateResp ContactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return nil, fmt.Errorf("decode update response: %w", err)
	}

	if updateResp.SyncToken != "" {
		c.syncToken = updateResp.SyncToken
	}
	if updateResp.PrefToken != "" {
		c.prefToken = updateResp.PrefToken
	}

	if len(updateResp.Contacts) == 0 {
		return nil, fmt.Errorf("update contact: no contact returned in response")
	}

	return &updateResp.Contacts[0], nil
}

func (c *Client) DeleteContact(contactID, etag string) error {
	if contactID == "" {
		return fmt.Errorf("delete contact: contactID is required")
	}
	if etag == "" {
		return fmt.Errorf("delete contact: etag is required")
	}

	if err := c.ensureContactsInit(); err != nil {
		return err
	}

	if c.contactsURL == "" || c.dsid == "" {
		return fmt.Errorf("contacts not available: missing contactsURL or dsid")
	}

	reqBody := struct {
		Contacts []struct {
			ContactID string `json:"contactId"`
			Etag      string `json:"etag"`
		} `json:"contacts"`
	}{
		Contacts: []struct {
			ContactID string `json:"contactId"`
			Etag      string `json:"etag"`
		}{
			{ContactID: contactID, Etag: etag},
		},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.buildContactsURLWithTokens("contacts/card/")+"&method=DELETE", bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("referer", "https://www.icloud.com/")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 450 {
		return ErrSessionExpired
	}

	if resp.StatusCode == 409 {
		return ErrContactEtagMismatch
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete contact failed with status %d: %s", resp.StatusCode, string(body))
	}

	var deleteResp ContactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
		return nil
	}

	if deleteResp.SyncToken != "" {
		c.syncToken = deleteResp.SyncToken
	}
	if deleteResp.PrefToken != "" {
		c.prefToken = deleteResp.PrefToken
	}

	return nil
}

// ensureContactsInit calls reqStartup to initialize sync tokens if they haven't been set yet.
func (c *Client) ensureContactsInit() error {
	if c.syncToken == "" || c.prefToken == "" {
		_, err := c.reqStartup()
		return err
	}
	return nil
}

func (c *Client) reqStartup() (*ContactsStartupResp, error) {
	if c.contactsURL == "" || c.dsid == "" {
		return nil, fmt.Errorf("contacts not available: missing contactsURL or dsid")
	}

	req, err := http.NewRequest(http.MethodGet, c.buildContactsURL("startup"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("referer", "https://www.icloud.com/")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 450 {
		return nil, ErrSessionExpired
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("contacts startup failed with status %d: %s", resp.StatusCode, string(body))
	}

	var startupResp ContactsStartupResp
	if err := json.NewDecoder(resp.Body).Decode(&startupResp); err != nil {
		return nil, fmt.Errorf("decode startup response: %w", err)
	}

	// Store tokens for subsequent mutation requests.
	c.syncToken = startupResp.SyncToken
	c.prefToken = startupResp.PrefToken

	return &startupResp, nil
}
