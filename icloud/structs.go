package icloud

type AuthInitResp struct {
	Iteration int    `json:"iteration"`
	Salt      string `json:"salt"`
	Protocol  string `json:"protocol"`
	B         string `json:"b"`
	C         string `json:"c"`
}

type AuthInitReq struct {
	A           string   `json:"a"`
	AccountName string   `json:"accountName"`
	Protocols   []string `json:"protocols"`
}

type AuthCompleteReq struct {
	AccountName string   `json:"accountName"`
	RememberMe  bool     `json:"rememberMe"`
	TrustTokens []string `json:"trustTokens"`
	M1          string   `json:"m1"`
	C           string   `json:"c"`
	M2          string   `json:"m2"`
}

type SecurityCode struct {
	Code string `json:"code"`
}

type PhoneNumberID struct {
	ID int `json:"id"`
}

type TwoFactorCodeFromPhoneRequest struct {
	SecurityCode *SecurityCode  `json:"securityCode,omitempty"`
	PhoneNumber  *PhoneNumberID `json:"phoneNumber,omitempty"`
	Mode         string         `json:"mode,omitempty"`
}

// ServiceError represents a Apple service error.
type ServiceError struct {
	Code    string `json:"code,omitempty"`
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

type HMEListResp struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Result    Result `json:"result"`
}

type Result struct {
	ForwardToEmails   []string   `json:"forwardToEmails"`
	HmeEmails         []HmeEmail `json:"hmeEmails"`
	SelectedForwardTo string     `json:"selectedForwardTo"`
}

type HmeEmail struct {
	Origin          string `json:"origin"`
	AnonymousID     string `json:"anonymousId"`
	Domain          string `json:"domain"`
	ForwardToEmail  string `json:"forwardToEmail"`
	Hme             string `json:"hme"`
	Label           string `json:"label"`
	Note            string `json:"note"`
	CreateTimestamp int64  `json:"createTimestamp"`
	IsActive        bool   `json:"isActive"`
	RecipientMailID string `json:"recipientMailId"`
	OriginAppName   string `json:"originAppName,omitempty"`
	AppBundleID     string `json:"appBundleId,omitempty"`
}

type MailInboxResp struct {
	TotalThreadsReturned int            `json:"totalThreadsReturned"`
	ThreadList           []Thread       `json:"threadList"`
	SessionHeaders       SessionHeaders `json:"sessionHeaders"`
	Events               []any          `json:"events"`
}

type Thread struct {
	JSONType           string   `json:"jsonType"`
	ThreadID           string   `json:"threadId"`
	Timestamp          int64    `json:"timestamp"`
	Count              int      `json:"count"`
	FolderMessageCount int      `json:"folderMessageCount"`
	Flags              []string `json:"flags"`
	Senders            []string `json:"senders"`
	Subject            string   `json:"subject"`
	Preview            string   `json:"preview"`
	Modseq             int64    `json:"modseq"`
}

type SessionHeaders struct {
	Folder       string `json:"folder"`
	Modseq       int64  `json:"modseq"`
	Threadmodseq int64  `json:"threadmodseq"`
	Condstore    int    `json:"condstore"`
	Qresync      int    `json:"qresync"`
	Threadmode   int    `json:"threadmode"`
}

type MessageMetadata struct {
	UID       string   `json:"uid"`
	Date      int64    `json:"date"`
	Size      int      `json:"size"`
	Folder    string   `json:"folder"`
	Modseq    int64    `json:"modseq"`
	Flags     []any    `json:"flags"`
	SentDate  string   `json:"sentDate"`
	Subject   string   `json:"subject"`
	From      []string `json:"from"`
	To        []string `json:"to"`
	Cc        []any    `json:"cc"`
	Bcc       []any    `json:"bcc"`
	Bimi      BIMI     `json:"bimi"`
	Smime     SMIME    `json:"smime"`
	MessageID string   `json:"messageId"`
	Parts     []Part   `json:"parts"`
}

type BIMI struct {
	Status string `json:"status"`
}

type SMIME struct {
	Status string `json:"status"`
}

type Part struct {
	JSONType    string `json:"jsonType"`
	ContentType string `json:"contentType"`
	Params      string `json:"params"`
	Encoding    string `json:"encoding"`
	Size        int    `json:"size"`
	PartID      string `json:"partId"`
	Lines       int    `json:"lines"`
	IsAttach    bool   `json:"isAttach"`
}

type MailDraftReq struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	From               string `json:"from"`
	To                 string `json:"to"`
	Subject            string `json:"subject"`
	TextBody           string `json:"textBody"`
	Body               string `json:"body"`
	Attachments        []any  `json:"attachments"`
	WebmailClientBuild string `json:"webmailClientBuild"`
}

type Message struct {
	GUID           string         `json:"guid"`
	LongHeader     string         `json:"longHeader"`
	To             []string       `json:"to"`
	From           []string       `json:"from"`
	Cc             []string       `json:"cc"`
	Bcc            []string       `json:"bcc"`
	ContentType    string         `json:"contentType"`
	Bimi           BIMI           `json:"bimi"`
	Smime          SMIME          `json:"smime"`
	Parts          []MsgPart      `json:"parts"`
	SessionHeaders SessionHeaders `json:"sessionHeaders"`
	Events         []any          `json:"events"`
}

type MsgPart struct {
	GUID    string `json:"guid"`
	Content string `json:"content"`
}
