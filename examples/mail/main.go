package main

import (
	"fmt"

	"github.com/Johnw7789/Go-iClient/icloud"
)

func main() {
	// Uncomment the example you want to run
	// RetrieveInbox()
	// RetrieveMessage()
	// DeleteMessage()
	// SendMail()
}

func promptOTP() (string, error) {
	var otp string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otp)
	return otp, nil
}

func RetrieveInbox() {
	client, err := icloud.NewClient("email", "password", true)
	if err != nil {
		panic(err)
	}

	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}

	mailResponse, err := client.RetrieveMailInbox(50, 0)
	if err != nil {
		panic(err)
	}

	for _, message := range mailResponse.ThreadList {
		fmt.Println(message.Senders)
		fmt.Println(message.Subject)
		fmt.Println(message.ThreadID)
		fmt.Println()
	}
}

func RetrieveMessage() {
	client, err := icloud.NewClient("email", "password", true)
	if err != nil {
		panic(err)
	}

	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}

	threadId := "threadId"

	mailMetadata, err := client.GetMessageMetadata(threadId)
	if err != nil {
		panic(err)
	}

	message, err := client.GetMessage(mailMetadata.UID)
	if err != nil {
		panic(err)
	}

	for _, part := range message.Parts {
		fmt.Println(part.Content)
	}
}

func DeleteMessage() {
	client, err := icloud.NewClient("email", "password", true)
	if err != nil {
		panic(err)
	}

	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}

	uid := "uid"

	success, err := client.DeleteMail(uid)
	if err != nil {
		panic(err)
	}

	fmt.Println("Email message deletion success:", success)
}

func SendMail() {
	client, err := icloud.NewClient("email", "password", true)
	if err != nil {
		panic(err)
	}

	if err := client.Login(promptOTP); err != nil {
		panic(err)
	}

	fromEmail := "test@icloud.com"
	toEmail := "test@icloud.com"
	subject := "Test Email"
	textBody := "This is a test email"
	body := "<html><body><h1>This is a test email</h1></body></html>"

	uid, err := client.DraftMail(fromEmail, toEmail, subject, textBody, body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Email draft uid:", uid)

	success, err := client.SendDraft(uid)
	if err != nil {
		panic(err)
	}

	fmt.Println("Email send success:", success)
}
