package main

import (
	"fmt"

	"github.com/Johnw7789/Go-iClient/icloud"
)

func main() {
	// client := newClient("username", "password")

	// Uncomment the example you want to run
	// RetrieveInbox(client)
	// RetrieveMessage(client, "thread-id")
	// DeleteMessage(client, "uid")
	// SendMail(client, "from@icloud.com", "to@example.com", "Subject", "plain text body", "<b>html body</b>")
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

func RetrieveInbox(client *icloud.Client) {
	mailResponse, err := client.RetrieveMailInbox(50, 0)
	if err != nil {
		panic(err)
	}

	for _, thread := range mailResponse.ThreadList {
		fmt.Println(thread.Senders)
		fmt.Println(thread.Subject)
		fmt.Println(thread.ThreadID)
		fmt.Println()
	}
}

func RetrieveMessage(client *icloud.Client, threadID string) {
	metadata, err := client.GetMessageMetadata(threadID)
	if err != nil {
		panic(err)
	}

	message, err := client.GetMessage(metadata.UID)
	if err != nil {
		panic(err)
	}

	for _, part := range message.Parts {
		fmt.Println(part.Content)
	}
}

func DeleteMessage(client *icloud.Client, uid string) {
	success, err := client.DeleteMail(uid)
	if err != nil {
		panic(err)
	}

	fmt.Println("Email deletion success:", success)
}

func SendMail(client *icloud.Client, from, to, subject, textBody, htmlBody string) {
	uid, err := client.DraftMail(from, to, subject, textBody, htmlBody)
	if err != nil {
		panic(err)
	}

	fmt.Println("Draft UID:", uid)

	success, err := client.SendDraft(uid)
	if err != nil {
		panic(err)
	}

	fmt.Println("Email send success:", success)
}
