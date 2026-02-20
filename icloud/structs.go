package icloud

import "encoding/json"

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

// ---- Find My ----

// FMDeviceFeatures maps feature flag codes (e.g. "SND", "LCK", "WIP") to their enabled state.
type FMDeviceFeatures map[string]bool

// FMAudioChannel represents a single audio channel on an AirPods device.
type FMAudioChannel struct {
	Name      string `json:"name"`
	Available int    `json:"available"`
	Playing   bool   `json:"playing"`
	Muted     bool   `json:"muted"`
}

// FMSndStatus is the status of a play-sound command on a device.
type FMSndStatus struct {
	AlertText           string `json:"alertText"`
	CancelButtonTitle   string `json:"cancelButtonTitle"`
	ContinueButtonTitle string `json:"continueButtonTitle"`
	CreateTimestamp     int64  `json:"createTimestamp"`
	StatusCode          string `json:"statusCode"`
	AlertTitle          string `json:"alertTitle"`
}

// FMMsgStatus is the status of a message/alert command on a device.
type FMMsgStatus struct {
	Vibrate         bool   `json:"vibrate"`
	Strobe          bool   `json:"strobe"`
	UserText        bool   `json:"userText"`
	PlaySound       bool   `json:"playSound"`
	CreateTimestamp int64  `json:"createTimestamp"`
	StatusCode      string `json:"statusCode"`
}

// FMLostDevice holds lost-mode details when a device is in lost mode.
type FMLostDevice struct {
	StopLostMode    bool   `json:"stopLostMode"`
	EmailUpdates    bool   `json:"emailUpdates"`
	UserText        bool   `json:"userText"`
	Sound           bool   `json:"sound"`
	OwnerNbr        string `json:"ownerNbr"`
	Text            string `json:"text"`
	Email           string `json:"email"`
	CreateTimestamp int64  `json:"createTimestamp"`
	StatusCode      string `json:"statusCode"`
}

// FMDevice represents a single device returned by the Find My service.
type FMDevice struct {
	ID                   string           `json:"id"`
	Name                 string           `json:"name"`
	DeviceDisplayName    string           `json:"deviceDisplayName"`
	DeviceClass          string           `json:"deviceClass"`
	DeviceModel          string           `json:"deviceModel"`
	RawDeviceModel       string           `json:"rawDeviceModel"`
	ModelDisplayName     string           `json:"modelDisplayName"`
	DeviceColor          string           `json:"deviceColor"`
	DeviceStatus         string           `json:"deviceStatus"`
	BatteryLevel         *float64         `json:"batteryLevel"`
	BatteryStatus        *string          `json:"batteryStatus"`
	LocationCapable      bool             `json:"locationCapable"`
	LocationEnabled      bool             `json:"locationEnabled"`
	IsLocating           bool             `json:"isLocating"`
	LostModeCapable      bool             `json:"lostModeCapable"`
	LostModeEnabled      bool             `json:"lostModeEnabled"`
	ActivationLocked     bool             `json:"activationLocked"`
	PasscodeLength       int              `json:"passcodeLength"`
	CanWipeAfterLock     bool             `json:"canWipeAfterLock"`
	WipeInProgress       bool             `json:"wipeInProgress"`
	IsMac                bool             `json:"isMac"`
	ThisDevice           bool             `json:"thisDevice"`
	FmlyShare            bool             `json:"fmlyShare"`
	LowPowerMode         bool             `json:"lowPowerMode"`
	DeviceWithYou        bool             `json:"deviceWithYou"`
	PendingRemove        bool             `json:"pendingRemove"`
	DarkWake             bool             `json:"darkWake"`
	MaxMsgChar           int              `json:"maxMsgChar"`
	BaUUID               string           `json:"baUUID"`
	DeviceDiscoveryId    string           `json:"deviceDiscoveryId"`
	CommandLookupId      string           `json:"commandLookupId"`
	LostTimestamp        string           `json:"lostTimestamp"`
	BrassStatus          string           `json:"brassStatus"`
	Features             FMDeviceFeatures `json:"features"`
	AudioChannels        []FMAudioChannel `json:"audioChannels"`
	Snd                  *FMSndStatus     `json:"snd"`
	Msg                  *FMMsgStatus     `json:"msg"`
	LostDevice           *FMLostDevice    `json:"lostDevice"`
	Location             interface{}      `json:"location"`
	TrackingInfo         interface{}      `json:"trackingInfo"`
	RemoteWipe           interface{}      `json:"remoteWipe"`
	RemoteLock           interface{}      `json:"remoteLock"`
	WipedTimestamp       interface{}      `json:"wipedTimestamp"`
	LockedTimestamp      interface{}      `json:"lockedTimestamp"`
	RepairStatus         interface{}      `json:"repairStatus"`
	EncodedDeviceId      interface{}      `json:"encodedDeviceId"`
	PrsId                interface{}      `json:"prsId"`
}

// FMDevicesResp is the response envelope from refreshClient and playSound.
// ServerContext is stored raw so it can be echoed back verbatim in subsequent requests.
type FMDevicesResp struct {
	UserInfo        json.RawMessage `json:"userInfo"`
	Alert           interface{}     `json:"alert"`
	ServerContext   json.RawMessage `json:"serverContext"`
	Content         []FMDevice      `json:"content"`
	UserPreferences json.RawMessage `json:"userPreferences"`
	StatusCode      string          `json:"statusCode"`
}

// ---- Find My end ----

// ---- Contacts ----

type ContactPhone struct {
	Label string `json:"label"`
	Field string `json:"field"`
}

type ContactEmail struct {
	Label string `json:"label"`
	Field string `json:"field"`
}

type ContactAddress struct {
	Label      string `json:"label,omitempty"`
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

type ContactURL struct {
	Label string `json:"label"`
	Field string `json:"field"`
}

type Contact struct {
	ContactID   string           `json:"contactId,omitempty"`
	Etag        string           `json:"etag,omitempty"`
	FirstName   string           `json:"firstName,omitempty"`
	LastName    string           `json:"lastName,omitempty"`
	MiddleName  string           `json:"middleName,omitempty"`
	NamePrefix  string           `json:"namePrefix,omitempty"`
	NameSuffix  string           `json:"nameSuffix,omitempty"`
	Nickname    string           `json:"nickname,omitempty"`
	Phones      []ContactPhone   `json:"phones,omitempty"`
	Emails      []ContactEmail   `json:"emails,omitempty"`
	Addresses   []ContactAddress `json:"addresses,omitempty"`
	URLs        []ContactURL     `json:"urls,omitempty"`
	Birthday    string           `json:"birthday,omitempty"`
	CompanyName string           `json:"companyName,omitempty"`
	JobTitle    string           `json:"jobTitle,omitempty"`
	Department  string           `json:"department,omitempty"`
	Notes       string           `json:"notes,omitempty"`
	IsCompany   bool             `json:"isCompany"`
}

type ContactsStartupResp struct {
	Contacts      []Contact `json:"contacts"`
	SyncToken     string    `json:"syncToken"`
	PrefToken     string    `json:"prefToken"`
	ContactsOrder []string  `json:"contactsOrder"`
	MeCardId      string    `json:"meCardId"`
}

type ContactsResponse struct {
	Contacts  []Contact `json:"contacts"`
	SyncToken string    `json:"syncToken,omitempty"`
	PrefToken string    `json:"prefToken,omitempty"`
}

// ---- Contacts end ----

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
