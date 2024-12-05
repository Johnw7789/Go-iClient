# Go iClient
Provides the ability to interact with icloud.com's web api. The authentication flow grants trust for future logins as well as auth for all icloud.com web applications.

An iCloud+ plan is required in order to use some features of the client, such as the Hide My Email service.

Please note that HME email generation is limited to 5 emails per hour, and up to 750 total emails per account, or if in the case of a family account, 750 per family member. 

## Progress / Roadmap

| **Module** | **Status** |
|:---:|:---:|
| Hide My Email |`✔`|
| Mail |`✔`|
| Find My |:hammer:| 
| Photos |:hammer:|
| iCloud Drive ||

## Installation
``go get github.com/Johnw7789/Go-iClient``

## Authentication

### Initializing an iCloud Client & Login
```
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

// * Wait for OTP, in this example we wait for a user to input it
var otpInput string
fmt.Print("Enter OTP: ")
fmt.Scanln(&otpInput)

// * Send the code back through the channel to the iClient
iclient.OtpChannel <- otpInput
```
## Usage

### Generating an HME email
```
label := "email 123"
note := "this email is a test email"

// * Generate a new HME email
emailAddress, err := iclient.ReserveHME(label, note)
if err != nil {
	log.Fatal(err)
}
```

### Retrieving all HME emails
```
// * Get all HME emails for the user
emails, err := iclient.RetrieveHMEList()
if err != nil {
	log.Fatal(err)
}

for _, email := range emails {
	fmt.Println(email.Hme)
}
```

### Deactivating an HME email
```
// * Anonymous Id for reactivation/deactivation/deletion can be retrieved from the HmeEmail struct as part of the HMEListResp struct
anonymousId := "id_here"

// * Deactivate the HME email
success, err := iclient.DeactivateHME(anonymousId)
if err != nil {
	log.Fatal(err)
}
```

### Reactivating an HME email
```
// * Reactivate the HME email
success, err = iclient.ReactivateHME(anonymousId)
if err != nil {
	log.Fatal(err)
}
```

### Deleting an HME email
```
// * In order to delete we must first deactive the HME email
success, err = iclient.DeactivateHME(anonymousId)
if err != nil {
	log.Fatal(err)
}

success, err = iclient.DeleteHME(anonymousId)
if err != nil {
	log.Fatal(err)
}
```

#### Retrieving the mail inbox
```
maxResults := 50
beforeTimestamp := 0 // if set to 0, it will be exluded from the query and set as a blank string

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
```
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
```
uid := "uid"

success, err := iclient.DeleteMail(uid)
if err != nil {
	log.Fatal(err)
}
```

### Sending an email 
```
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
