package icloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	http "github.com/bogdanfinn/fhttp"
	"github.com/tidwall/gjson"
)

// * RetrieveMailInbox() retrieves the user's inbox from their iCloud account, maxResults is the maximum number of emails to retrieve, beforeTs is the timestamp to retrieve emails before
func (c *Client) RetrieveMailInbox(maxResults, beforeTs int) (MailInboxResp, error) {
	body, err := c.reqRetrieveMailInbox(maxResults, beforeTs)
	if err != nil {
		return MailInboxResp{}, err
	}

	var mailInboxResp MailInboxResp
	err = json.Unmarshal([]byte(body), &mailInboxResp)
	if err != nil {
		return MailInboxResp{}, err
	}

	return mailInboxResp, nil
}

// * GetMessage() retrieves the message metadata from the user's inbox using the threadId
func (c *Client) GetMessageMetadata(threadId string) (MessageMetadata, error) {
	body, err := c.reqGetMessageMetadata(threadId)
	if err != nil {
		return MessageMetadata{}, err
	}

	msgMdArr := gjson.Get(body, "messageMetadataList").Array()
	if len(msgMdArr) == 0 {
		return MessageMetadata{}, errors.New("no message metadata found")
	}

	for _, msgMd := range msgMdArr {
		if msgMd.Get("folder").String() == "INBOX" {
			msgMdStr := msgMd.String()

			var messageMetadata MessageMetadata
			err = json.Unmarshal([]byte(msgMdStr), &messageMetadata)
			if err != nil {
				return MessageMetadata{}, err
			}

			return messageMetadata, nil
		}
	}

	return MessageMetadata{}, errors.New("failed to get message metadata")
}

// * GetMessage() retrieves the message from the user's inbox using the uid, and returns the full body html of the email
func (c *Client) GetMessage(uid string) (Message, error) {
	body, err := c.reqGetMessage(uid)
	if err != nil {
		return Message{}, err
	}

	var message Message
	err = json.Unmarshal([]byte(body), &message)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}

// * DeleteMail() deletes an email from the user's inbox using the uid
func (c *Client) DeleteMail(uid string) (bool, error) {
	body, err := c.reqMailDelete(uid)
	if err != nil {
		return false, err
	}

	// todo: is this fine for comfirming deletion? since this will still appear in resp after already being deleted before
	uidDeleted := (gjson.Get(body, "result.deletedUids").String() == uid)

	return uidDeleted, nil
}

// * DraftMail() drafts and saves an email to the user's draft folder. fromName is the name of the sender, fromEmail is the email of the sender, toName is the name of the recipient, toEmail is the email of the recipient, subject is the subject of the email, textBody is the text body of the email, body is the full html body of the email
func (c *Client) DraftMail(fromEmail, toEmail, subject, textBody, body string) (string, error) {
	body, err := c.reqMailDraft(fromEmail, toEmail, subject, textBody, body)
	if err != nil {
		return "", err
	}

	uid := gjson.Get(body, "result.uid").String()

	// since gjosn.Get returns an empty string if not found, we can use this to check if the email was successfully drafted
	if uid == "" {
		return "", errors.New("failed to draft email")
	}

	return uid, nil
}

// * SendDraft() sends an email from the user's draft folder using the uid
func (c *Client) SendDraft(uid string) (bool, error) {
	resp, err := c.reqSendDraft(uid)
	if err != nil {
		return false, err
	}

	// * 404 is the default response for failure/draft not found
	return resp.StatusCode == 200, nil
}

func (c *Client) reqRetrieveMailInbox(maxResults, beforeTs int) (string, error) {
	beforeTsStr := fmt.Sprintf("%d", beforeTs)

	if beforeTs == 0 {
		beforeTsStr = ""
	}

	body := bytes.NewReader([]byte(fmt.Sprintf(`{"responseType":"THREAD_DIGEST","includeFolderStatus":false,"maxResults":%d,"before":"%s","sessionHeaders":{"folder":"INBOX","condstore":1,"qresync":1,"threadmode":1}}`, maxResults, beforeTsStr)))

	req, err := http.NewRequest(http.MethodPost, endpoints[mailInbox], body)
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

func (c *Client) reqGetMessageMetadata(threadId string) (string, error) {
	body := bytes.NewReader([]byte(`{"threadId":"` + threadId + `","includeLabelIds":false,"sessionHeaders":{"folder":"INBOX","condstore":1,"qresync":1,"threadmode":1}}`))

	req, err := http.NewRequest(http.MethodPost, endpoints[mailMetadataGet], body)
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

func (c *Client) reqGetMessage(uid string) (string, error) {
	body := bytes.NewReader([]byte(`{"uid":"` + uid + `","parts":["2.1"],"dontMarkAsRead":true,"sessionHeaders":{"folder":"INBOX","condstore":1,"qresync":1,"threadmode":1}}`))

	req, err := http.NewRequest(http.MethodPost, endpoints[mailGet], body)
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

func (c *Client) reqMailDelete(uid string) (string, error) {
	body := bytes.NewReader([]byte(`{"jsonrpc":"2.0","method":"delete","params":{"folder":"folder:INBOX","uids":["` + uid + `"],"rollbackslot":"0.0"}}`))

	req, err := http.NewRequest(http.MethodPost, endpoints[mailDelete], body)
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

// todo: support attachments
func (c *Client) reqMailDraft(fromEmail, toEmail, subject, textBody, body string) (string, error) {
	payload := MailDraftReq{
		Jsonrpc: "2.0",
		Method:  "saveDraft",
		Params: Params{
			From:               fromEmail,
			To:                 toEmail,
			Subject:            subject,
			TextBody:           textBody,
			Body:               body,
			Attachments:        []any{},
			WebmailClientBuild: "current",
		},
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // Disable HTML escaping
	if err := encoder.Encode(payload); err != nil {
		return "", err
	}

	jsonBytes := buf.Bytes()

	req, err := http.NewRequest(http.MethodPost, endpoints[mailDraft], bytes.NewReader(jsonBytes))
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

func (c *Client) reqSendDraft(uid string) (*http.Response, error) {
	body := bytes.NewReader([]byte(`{"messageGuid":"Drafts/` + uid + `","sessionHeaders":{"folder":null,"modseq":null,"threadmodseq":null,"condstore":1,"qresync":1,"threadmode":1}}`))

	req, err := http.NewRequest(http.MethodPost, endpoints[mailSend], body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(HdrContentType, "application/json")
	req.Header.Set("origin", "https://www.icloud.com")
	req.Header.Set("accept", "*/*")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
