package icloud

import (
	"errors"

	http "github.com/bogdanfinn/fhttp"
)

const (
	ProtocolVersion = "QH5B2"
	OAuthClientID   = "d39ba9916b7251055b22c7f910e2ea796ee65e98b2ddecea8f5dde8d9d1a815d"
)

const (
	HdrContentType       = "Content-Type"
	HdrXRequestedWith    = "X-Requested-With"
	HdrXAppleWidgetKey   = "X-Apple-Widget-Key"
	HdrAccept            = "Accept"
	HdrXAppleIDSessionID = "X-Apple-ID-Session-Id"
	HdrScnt              = "scnt"
)

var (
	ErrIncorrectUsernamePassword = errors.New("incorrect username or password")
	ErrRequiredPrivacyAck        = errors.New("sign in to https://appleid.apple.com and acknowledge the Apple ID and Privacy agreement")
	ErrUnexpectedSigninResponse  = errors.New("unexpected sign in response")
	ErrNotImplemented            = errors.New("not implemented")
	ErrSeverErrorOrInvalidCreds  = errors.New("apple server error or invalid credentials")
	ErrFindMySessionExpired      = errors.New("find my session expired: please call Login() again")
	ErrSessionExpired            = errors.New("icloud session expired: please call Login() again")
	ErrContactEtagMismatch = errors.New("contact etag mismatch: contact was modified")
)

type endpoint uint8

const (
	iTunesConnect endpoint = 1 + iota
	authStart
	authFederate
	authInit
	authComplete
	authOptions
	submitSecurityCode
	trust
	authWeb
	authValidate

	hmeList
	hmeGen
	hmeReserve
	hmeDeactivate
	hmeReactivate
	hmeDelete

	mailInbox
	mailMetadataGet
	mailGet
	mailDelete
	mailDraft
	mailSend
)

var endpoints = map[endpoint]string{
	authStart:          "https://idmsa.apple.com/appleauth/auth/authorize/signin?frame_id=auth-%s&language=en_US&skVersion=7&iframeId=auth-%s&client_id=%s&redirect_uri=https://www.icloud.com&response_type=code&response_mode=web_message&state=auth-%s&authVersion=latest",
	authFederate:       "https://idmsa.apple.com/appleauth/auth/federate?isRememberMeEnabled=true",
	authInit:           "https://idmsa.apple.com/appleauth/auth/signin/init",
	authComplete:       "https://idmsa.apple.com/appleauth/auth/signin/complete?isRememberMeEnabled=true",
	authOptions:        "https://idmsa.apple.com/appleauth/auth",
	submitSecurityCode: "https://idmsa.apple.com/appleauth/auth/verify/%s/securitycode", // code type, typically trusteddevice
	trust:              "https://idmsa.apple.com/appleauth/auth/2sv/trust",
	authWeb:            "https://setup.icloud.com/setup/ws/1/accountLogin",
	authValidate:       "https://setup.icloud.com/setup/ws/1/validate?clientBuildNumber=2602Build17&clientMasteringNumber=2602Build17&clientId=%s&dsid=%s",

	hmeList:       "https://p52-maildomainws.icloud.com/v2/hme/list?clientBuildNumber=2426Hotfix51&clientMasteringNumber=2426Hotfix51",
	hmeGen:        "https://p52-maildomainws.icloud.com/v1/hme/generate?clientBuildNumber=2415Project29&clientMasteringNumber=2415B20",
	hmeReserve:    "https://p52-maildomainws.icloud.com/v1/hme/reserve?clientBuildNumber=2415Project29&clientMasteringNumber=2415B20",
	hmeDeactivate: "https://p52-maildomainws.icloud.com/v1/hme/deactivate?clientBuildNumber=2426Hotfix10&clientMasteringNumber=2426Hotfix10",
	hmeReactivate: "https://p52-maildomainws.icloud.com/v1/hme/reactivate?clientBuildNumber=2426Hotfix10&clientMasteringNumber=2426Hotfix10",
	hmeDelete:     "https://p52-maildomainws.icloud.com/v1/hme/delete?clientBuildNumber=2426Hotfix10&clientMasteringNumber=2426Hotfix10",

	mailInbox:       "https://p52-mccgateway.icloud.com/mailws2/v1/thread/search?clientBuildNumber=2426Hotfix40&clientMasteringNumber=2426Hotfix40",
	mailMetadataGet: "https://p52-mccgateway.icloud.com/mailws2/v1/thread/get?clientBuildNumber=2426Hotfix40&clientMasteringNumber=2426Hotfix40",
	mailGet:         "https://p52-mccgateway.icloud.com/mailws2/v1/message/get?clientBuildNumber=2426Hotfix40&clientMasteringNumber=2426Hotfix40",
	mailDelete:      "https://p52-mailws.icloud.com/wm/message?clientBuildNumber=2426Hotfix40&clientMasteringNumber=2426Hotfix40",
	mailDraft:       "https://p52-mailws.icloud.com/wm/message?clientBuildNumber=2426Hotfix40&clientMasteringNumber=2426Hotfix40",
	mailSend:        "https://p52-mccgateway.icloud.com/mailws2/v1/draft/send?clientBuildNumber=2426Hotfix40&clientMasteringNumber=2426Hotfix40",
}

// updateRequestHeaders updates required request headers.
func (c *Client) updateRequestHeaders(header http.Header) http.Header {
	if c.scnt != "" {
		header.Set(HdrScnt, c.scnt)
	}
	if c.sessionID != "" {
		header.Set(HdrXAppleIDSessionID, c.sessionID)
	}

	header.Set(HdrXRequestedWith, "XMLHttpRequest")
	header.Set(HdrContentType, "application/json")
	header.Set(HdrAccept, "application/json")
	header.Set("referer", "https://idmsa.apple.com/")
	header.Set("origin", "https://idmsa.apple.com")
	header.Set("X-Apple-Widget-Key", c.serviceKey)
	header.Set("X-Requested-With", "XMLHttpRequest")
	header.Set("X-Apple-I-Require-UE", "true")
	header.Set("X-Apple-Auth-Attributes", c.authAttr)
	header.Set("X-Apple-Widget-Key", c.clientId)
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	header.Set("X-Apple-Mandate-Security-Upgrade", "0")
	header.Set("X-Apple-Oauth-Client-Id", c.clientId)
	header.Set("X-Apple-I-FD-Client-Info", `{"U":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36","L":"en-US","Z":"GMT-04:00","V":"1.1","F":".ta44j1e3NlY5BNlY5BSs5uQ32SCVgdI.AqWJ4EKKw0fVD_DJhCizgzH_y3EjNklY_ia4WFL264HRe4FSr_JzC1zJ6rgNNlY5BNp55BNlan0Os5Apw.BS1"}`)
	header.Set("X-Apple-Oauth-Client-Type", `firstPartyAuth`)
	header.Set("X-Apple-Oauth-Redirect-URI", `https://www.icloud.com`)
	header.Set("X-Apple-Oauth-Require-Grant-Code", `true`)
	header.Set("X-Apple-Oauth-Response-Mode", `web_message`)
	header.Set("X-Apple-Oauth-Response-Type", `code`)
	header.Set("X-Apple-Oauth-State", `auth-`+c.frameId)
	header.Set("X-Apple-Offer-Security-Upgrade", `1`)
	header.Set("X-Apple-Frame-Id", `auth-`+c.frameId)

	return header
}
