# Go iClient
Provides the ability to interact with icloud.com's web api. The authentication flow grants trust for future logins as well as auth for all icloud.com web applications.

An iCloud+ plan is required in order to use some features of the client, such as the Hide My Email service.

Please note that HME email generation is limited to 5 emails per hour, and up to 750 total emails per account, or if in the case of a family account, 750 per family member. 

## Progress / Roadmap

| **Module** | **Status** |
|:---:|:---:|
| Hide My Email |`✔`|
| Mail |`✔`|
| Find My |`✔`|
| Contacts |`✔`|
| Photos |:hammer:|
| iCloud Drive |:hammer:|

## Installation
``go get github.com/Johnw7789/Go-iClient``

## Authentication

### Initializing an iCloud Client & Login
In this example we are waiting for a OTP from user input from the console. However, this could be implemented in other ways. For example, if using this package as part of a project with a UI, a popup input modal could be opened to wait for the OTP, and then send the OTP back through the channel to the thread. 
```Go
// * Create a new iClient with account username and password, do not sniff with local proxy
iclient, err := icloud.NewClient("username", "password", false)
if err != nil {
	log.Fatal(err)
}

go func() {
	// * Create iCloud session and authenticate to use HME service
	err := iclient.Login()
	if err != nil {
		log.Fatal(err)
	}

	// do stuff after auth here e.g. gen hme, fetch mail etc
}()

var otpInput string
fmt.Print("Enter OTP: ")
fmt.Scanln(&otpInput)

iclient.OtpChannel <- otpInput
```
## Usage

### Generating an HME email
```Go
label := "email 123"
note := "this email is a test email"

// * Generate a new HME email
emailAddress, err := iclient.ReserveHME(label, note)
if err != nil {
	log.Fatal(err)
}
```

### Retrieving all HME emails
```Go
emails, err := iclient.RetrieveHMEList()
if err != nil {
	log.Fatal(err)
}

for _, email := range emails {
	fmt.Println(email.Hme)
}
```

### Deactivating an HME email
The anonymous ID is used for reactivation/deactivation/deletion and can be retrieved from the HmeEmail struct as part of the HMEListResp struct.
```Go
anonymousId := "anonymousId"

success, err := iclient.DeactivateHME(anonymousId)
if err != nil {
	log.Fatal(err)
}
```

### Reactivating an HME email
```Go
success, err = iclient.ReactivateHME(anonymousId)
if err != nil {
	log.Fatal(err)
}
```

### Deleting an HME email
In order to delete an email it first must be deactivated. 
```Go
success, err = iclient.DeactivateHME(anonymousId)
if err != nil {
	log.Fatal(err)
}

success, err = iclient.DeleteHME(anonymousId)
if err != nil {
	log.Fatal(err)
}
```

### Retrieving the mail inbox
```Go
maxResults := 50
beforeTimestamp := 0 // if set to 0, it will not be used in the query and instead set as a blank string

mailResponse, err := iclient.RetrieveMailInbox(maxResults, beforeTimestamp)
if err != nil {
	log.Fatal(err)
}

for _, message := range mailResponse.ThreadList {
	fmt.Println(message.Senders)
	fmt.Println(message.Subject)
	fmt.Println(message.ThreadID)
	fmt.Println()
}
```

### Retrieving an individual message
```Go
threadId := "threadId"

mailMetadata, err := iclient.GetMessageMetadata(threadId)
if err != nil {
	log.Fatal(err)
}

message, err := iclient.GetMessage(mailMetadata.UID)
if err != nil {
	log.Fatal(err)
}

for _, part := range message.Parts {
	// * Print the content of the email, or the body html
	fmt.Println(part.Content)
}
```

### Deleting an email 
The UID can only be obtained from the mail metadata, which is why you must get the message metadata first. 
```Go
threadId := "threadId"

mailMetadata, err := iclient.GetMessageMetadata(threadId)
if err != nil {
	log.Fatal(err)
}

success, err := iclient.DeleteMail(mailMetadata.UID)
if err != nil {
	log.Fatal(err)
}
```

### Sending an email
```Go
fromEmail := "test@icloud.com"
toEmail := "test@icloud.com"
subject := "Test Email"
textBody := "This is a test email"
body := "<html><body><h1>This is a test email</h1></body></html>"

uid, err := iclient.DraftMail(fromEmail, toEmail, subject, textBody, body)
if err != nil {
	log.Fatal(err)
}

// * Send the email
success, err := iclient.SendDraft(uid)
if err != nil {
	log.Fatal(err)
}
```

### Fetching all Find My devices
Includes all devices on the account and family sharing members.
```Go
devices, err := iclient.GetDevices()
if err != nil {
	log.Fatal(err)
}

for _, d := range devices {
	fmt.Printf("[%s] %s — %s\n", d.DeviceClass, d.Name, d.DeviceDisplayName)
}
```

### Fetching a single device
```Go
device, err := iclient.GetDevice(deviceID)
if err != nil {
	log.Fatal(err)
}
```

### Playing a sound on a device
Pass `nil` for the channels argument on iPhone, Watch, and Mac. For AirPods and other multi-channel accessories pass the channel names.
```Go
// iPhone, Watch, Mac
updated, err := iclient.PlaySound(deviceID, nil)

// AirPods
updated, err := iclient.PlaySound(deviceID, []string{"left", "right"})
```

### Keeping the session alive
`KeepAlive` blocks and calls the iCloud session validate endpoint on the given interval. Run it in a goroutine. It returns `ErrSessionExpired` if the session dies, or `ctx.Err()` if cancelled.
```Go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
	if err := iclient.KeepAlive(ctx, 5*time.Minute); err != nil {
		log.Println("session expired:", err)
	}
}()
```

### Fetching all contacts
GetContacts automatically initializes the contacts session on first call.
```Go
contacts, err := iclient.GetContacts()
if err != nil {
	log.Fatal(err)
}

for _, c := range contacts {
	fmt.Printf("%s %s\n", c.FirstName, c.LastName)
}
```

### Fetching a single contact
```Go
contact, err := iclient.GetContact(contactID)
if err != nil {
	log.Fatal(err)
}
```

### Creating a contact
```Go
newContact := icloud.Contact{
	FirstName: "John",
	LastName:  "Doe",
	Emails: []icloud.ContactEmail{
		{Label: "HOME", Field: "john@example.com"},
	},
	Phones: []icloud.ContactPhone{
		{Label: "MOBILE", Field: "+1234567890"},
	},
}

created, err := iclient.CreateContact(newContact)
if err != nil {
	log.Fatal(err)
}
```

### Updating a contact
Requires the contact's current etag to prevent concurrent modification conflicts.
```Go
contact.FirstName = "Jane"
updated, err := iclient.UpdateContact(contact)
if err != nil {
	log.Fatal(err)
}
```

### Deleting a contact
```Go
err := iclient.DeleteContact(contactID, etag)
if err != nil {
	log.Fatal(err)
}
```
