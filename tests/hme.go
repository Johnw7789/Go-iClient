package tests

import (
	"fmt"
	"strings"

	"github.com/Johnw7789/Go-iClient/icloud"
)

// * Anonymous Id for reactivation/deactivation/deletion can be retrieved from the HmeEmail struct as part of the HMEListResp struct

func TestReserveHME() {
	iclient, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		label := "email 123"
		note := "this email is a test email"

		// get hme from iclient.GenerateHME
		hme, err := iclient.ReserveHME(label, note)
		if err != nil {
			panic(err)
		}

		if !strings.Contains(hme, "@") {
			panic("Invalid email address")
		}
	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}

func TestRetrieveHMEList() {
	iclient, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		// get hme from iclient.GenerateHME
		emails, err := iclient.RetrieveHMEList()
		if err != nil {
			panic(err)
		}

		for _, email := range emails {
			fmt.Println(email.Hme)
		}
	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}

func TestDeactivateHME() {
	iclient, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		anonymousId := ""

		success, err := iclient.DeactivateHME(anonymousId)
		if err != nil {
			panic(err)
		}

		fmt.Println("HME deactivation success:", success)
	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}

func TestReactivateHME() {
	iclient, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		anonymousId := ""

		success, err := iclient.ReactivateHME(anonymousId)
		if err != nil {
			panic(err)
		}

		fmt.Println("HME reactivation success:", success)
	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}

func TestDeleteHME() {
	iclient, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		anonymousId := ""

		success, err := iclient.DeleteHME(anonymousId)
		if err != nil {
			panic(err)
		}

		fmt.Println("HME deletion success:", success)
	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}
