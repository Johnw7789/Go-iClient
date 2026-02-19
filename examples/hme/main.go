package main

import (
	"fmt"
	"strings"

	"github.com/Johnw7789/Go-iClient/icloud"
)

func main() {
	// client := newClient("username", "password")

	// Uncomment the example you want to run
	// RetrieveHMEList(client)
	// ReserveHME(client, "My Label", "optional note")
	// DeactivateHME(client, "anonymous-id")
	// ReactivateHME(client, "anonymous-id")
	// DeleteHME(client, "anonymous-id")
}

func promptOTP() (string, error) {
	var otp string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otp)
	return otp, nil
}

func newClient(username, password string) *icloud.Client {
	client, err := icloud.NewClient(username, password, false)
	if err != nil {
		panic(err)
	}
	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}
	return client
}

func RetrieveHMEList(client *icloud.Client) {
	emails, err := client.RetrieveHMEList()
	if err != nil {
		panic(err)
	}

	for _, email := range emails {
		fmt.Println(email.Hme)
	}
}

func ReserveHME(client *icloud.Client, label, note string) {
	hme, err := client.ReserveHME(label, note)
	if err != nil {
		panic(err)
	}

	if !strings.Contains(hme, "@") {
		panic("invalid email address returned")
	}

	fmt.Println("Reserved HME:", hme)
}

func DeactivateHME(client *icloud.Client, anonymousID string) {
	success, err := client.DeactivateHME(anonymousID)
	if err != nil {
		panic(err)
	}

	fmt.Println("HME deactivation success:", success)
}

func ReactivateHME(client *icloud.Client, anonymousID string) {
	success, err := client.ReactivateHME(anonymousID)
	if err != nil {
		panic(err)
	}

	fmt.Println("HME reactivation success:", success)
}

func DeleteHME(client *icloud.Client, anonymousID string) {
	success, err := client.DeleteHME(anonymousID)
	if err != nil {
		panic(err)
	}

	fmt.Println("HME deletion success:", success)
}
