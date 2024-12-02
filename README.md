# iCloud HME Gen
Provides the ability to generate iCloud emails using iCloud's Hide My Email functionality.

An iCloud+ plan is required in order to use the HME service.

Please note that email generation is limited to 5 emails per hour, and up to 750 total emails per account, or if in the case of a family account, 750 per family member. 

## Installation
``go get github.com/Johnw7789/iCloud-HME-Gen``

## Usage
```
// * Create a new iClient with account username and password
iclient, err := icloud.NewClient("user", "pass", false)
if err != nil {
	panic(err)
}

go func() {
	// * Create iCloud session and authenticate to use HME service
	err := iclient.Login()
	if err != nil {
		panic(err)
	}

	// * The prefix for the name of the email, set to whatever is desired
	identifierPrefix := "goicloud"

	// * Request a new email
	hme, err := iclient.GenerateHME(identifierPrefix)
	if err != nil {
		panic(err)
	}

	fmt.Println("Reserve HME success, hme: ", hme)
}()

// * Wait for user to input OTP from device
var otpInput string
fmt.Print("Enter OTP: ")
fmt.Scanln(&otpInput)

iclient.OtpChannel <- otpInput
```
