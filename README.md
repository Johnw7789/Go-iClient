# Go iClient
Provides the ability to interact with icloud.com's web api. The authentication flow grants trust for future logins as well as auth for all icloud.com web applications.

An iCloud+ plan is required in order to use some features of the client, such as the Hide My Email service.

Please note that email generation is limited to 5 emails per hour, and up to 750 total emails per account, or if in the case of a family account, 750 per family member. 

## Progress / Roadmap

| **Module** | **Status** |
|:---:|:---:|
| Hide My Email |`✔`|
| Mail |:hammer:|
| Find My |:hammer:| 
| Photos ||
| iCLoud Drive ||

## Installation
``go get github.com/Johnw7789/Go-iClient``

## HME Usage
```
// * Create a new iClient with account username and password, do not sniff with local proxy
iclient, err := icloud.NewClient("username", "password", false)
if err != nil {
	panic(err)
}

go func() {
	// * Create iCloud session and authenticate to use HME service
	err := iclient.Login()
	if err != nil {
		panic(err)
	}

	label := "email 123"
	note := "this email is a test email"

	// * Generate a new HME email
	hme, err := iclient.ReserveHME(label, note)
	if err != nil {
		panic(err)
	}

	fmt.Println("Reserve HME success, hme: ", hme)

	// * Get all HME emails for the user
	emails, err := iclient.RetrieveHMEList()
	if err != nil {
		panic(err)
	}

	for _, email := range emails {
		fmt.Println(email.Hme)
	}

	// * Anonymous Id for reactivation/deactivation/deletion can be retrieved from the HmeEmail struct as part of the HMEListResp struct
	anonymousId := "id_here"

	// * Deactivate the HME email
	success, err := iclient.DeactivateHME(anonymousId)
	if err != nil {
		panic(err)
	}

	fmt.Println("HME deactivation success:", success)

	// * Reactivate the HME email
	success, err = iclient.ReactivateHME(anonymousId)
	if err != nil {
		panic(err)
	}

	fmt.Println("HME reactivation success:", success)

	// * In order to delete we must first deactive the HME email
	success, err = iclient.DeactivateHME(anonymousId)
	if err != nil {
		panic(err)
	}

	success, err = iclient.DeleteHME(anonymousId)
	if err != nil {
		panic(err)
	}

	fmt.Println("HME deletion success:", success)
}()

// * Wait for OTP, in this example we wait for a user to input it
var otpInput string
fmt.Print("Enter OTP: ")
fmt.Scanln(&otpInput)

// * Send the code back through the channel to the iClient
iclient.OtpChannel <- otpInput
```
