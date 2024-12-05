package tests

import (
	"fmt"

	"github.com/Johnw7789/Go-iClient/icloud"
)

func TestRetrieveInbox() {
	iclient, err := icloud.NewClient("email", "password", true)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		mailResponse, err := iclient.RetrieveMailInbox(50, 0)
		if err != nil {
			panic(err)
		}

		for _, message := range mailResponse.ThreadList {
			fmt.Println(message.Senders)
			fmt.Println(message.Subject)
			fmt.Println(message.ThreadID)
			fmt.Println()
		}
	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}

func TestRetrieveMessage() {
	iclient, err := icloud.NewClient("email", "password", true)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		threadId := "threadId"

		mailMetadata, err := iclient.GetMessageMetadata(threadId)
		if err != nil {
			panic(err)
		}

		message, err := iclient.GetMessage(mailMetadata.UID)
		if err != nil {
			panic(err)
		}

		for _, part := range message.Parts {
			// * Print the content of the email, or the body html
			fmt.Println(part.Content)
		}

	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}

func TestDeleteMessage() {
	iclient, err := icloud.NewClient("email", "password", true)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		uid := "uid"

		success, err := iclient.DeleteMail(uid)
		if err != nil {
			panic(err)
		}

		fmt.Println("Email message deletion success:", success)
	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}

func TestSendMail() {
	iclient, err := icloud.NewClient("email", "password", true)
	if err != nil {
		panic(err)
	}

	go func() {
		err := iclient.Login()
		if err != nil {
			panic(err)
		}

		fromEmail := "test@icloud.com"
		toEmail := "test@icloud.com"
		subject := "Test Email"
		textBody := "This is a test email"
		body := "<html><body><h1>This is a test email</h1></body></html>"

		uid, err := iclient.DraftMail(fromEmail, toEmail, subject, textBody, body)
		if err != nil {
			panic(err)
		}

		fmt.Println("Email draft uid:", uid)

		// * Send the email
		success, err := iclient.SendDraft(uid)
		if err != nil {
			panic(err)
		}

		fmt.Println("Email send success:", success)
	}()

	// get otp from user input fmt.Scanln
	var otpInput string
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otpInput)

	iclient.OtpChannel <- otpInput

	select {}
}
