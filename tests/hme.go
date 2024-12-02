package tests

import (
	"fmt"

	"github.com/Johnw7789/iCloud-HME-Gen/icloud"
)

func TestHME() {
	iclient, err := icloud.NewClient("user", "pass", false)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		// * The prefix for the name of the email
		identifierPrefix := "goicloud"

		// get hme from iclient.GenerateHME
		hme, err := iclient.GenerateHME(identifierPrefix)
		if err != nil {
			panic(err)
		}

		fmt.Println("Reserve HME success, hme: ", hme)
	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}
