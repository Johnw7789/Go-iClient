package icloud

type authCompleteReq struct {
	AccountName string   `json:"accountName"`
	RememberMe  bool     `json:"rememberMe"`
	TrustTokens []string `json:"trustTokens"`
	M1          string   `json:"m1"`
	C           string   `json:"c"`
	M2          string   `json:"m2"`
}

type authInitResp struct {
	Iteration int    `json:"iteration"`
	Salt      string `json:"salt"`
	Protocol  string `json:"protocol"`
	B         string `json:"b"`
	C         string `json:"c"`
}

type authInitReq struct {
	A           string   `json:"a"`
	AccountName string   `json:"accountName"`
	Protocols   []string `json:"protocols"`
}

type securityCode struct {
	Code string `json:"code"`
}

type phoneNumberID struct {
	ID int `json:"id"`
}

type twoFactorCodeFromPhoneRequest struct {
	SecurityCode *securityCode  `json:"securityCode,omitempty"`
	PhoneNumber  *phoneNumberID `json:"phoneNumber,omitempty"`
	Mode         string         `json:"mode,omitempty"`
}
