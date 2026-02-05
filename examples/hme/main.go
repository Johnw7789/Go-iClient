package main

import (
	"fmt"
	"strings"

	"github.com/Johnw7789/Go-iClient/icloud"
)

func main() {
	// Uncomment the example you want to run
	// ReserveHME()
	// RetrieveHMEList()
	// DeactivateHME()
	// ReactivateHME()
	// DeleteHME()
}

func promptOTP() (string, error) {
	var otp string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otp)
	return otp, nil
}

func ReserveHME() {
	client, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}

	label := "email 123"
	note := "this email is a test email"

	hme, err := client.ReserveHME(label, note)
	if err != nil {
		panic(err)
	}

	if !strings.Contains(hme, "@") {
		panic("Invalid email address")
	}

	fmt.Println("Reserved HME:", hme)
}

func RetrieveHMEList() {
	client, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}

	emails, err := client.RetrieveHMEList()
	if err != nil {
		panic(err)
	}

	for _, email := range emails {
		fmt.Println(email.Hme)
	}
}

func DeactivateHME() {
	client, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}

	anonymousId := ""

	success, err := client.DeactivateHME(anonymousId)
	if err != nil {
		panic(err)
	}

	fmt.Println("HME deactivation success:", success)
}

func ReactivateHME() {
	client, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}

	anonymousId := ""

	success, err := client.ReactivateHME(anonymousId)
	if err != nil {
		panic(err)
	}

	fmt.Println("HME reactivation success:", success)
}

func DeleteHME() {
	client, err := icloud.NewClient("username", "password", true)
	if err != nil {
		panic(err)
	}

	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}

	anonymousId := ""

	success, err := client.DeleteHME(anonymousId)
	if err != nil {
		panic(err)
	}

	fmt.Println("HME deletion success:", success)
}
