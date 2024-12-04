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
	Success   bool      `json:"success"`
	Timestamp int       `json:"timestamp"`
	Result    Result    `json:"result"`
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
