package main

import (
	"fmt"

	"github.com/Johnw7789/Go-iClient/icloud"
)

func main() {
	// client := newClient("username", "password")

	// Uncomment the example you want to run
	// GetAllContacts(client)
	// CreateNewContact(client)
	// UpdateExistingContact(client, "contact-id", "current-etag")
	// DeleteExistingContact(client, "contact-id", "current-etag")
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

func GetAllContacts(client *icloud.Client) {
	contacts, err := client.GetContacts()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d contacts\n", len(contacts))

	for _, c := range contacts {
		fmt.Printf("- %s %s (%s), Etag: %s\n", c.FirstName, c.LastName, c.ContactID, c.Etag)
	}
}

func CreateNewContact(client *icloud.Client) {
	contact := icloud.Contact{
		FirstName: "John",
		LastName:  "Doe",
		Emails: []icloud.ContactEmail{
			{Label: "HOME", Field: "john.doe@example.com"},
		},
		Phones: []icloud.ContactPhone{
			{Label: "MOBILE", Field: "+1234567890"},
		},
	}

	created, err := client.CreateContact(contact)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created contact: %s (ID: %s, Etag: %s)\n", created.FirstName, created.ContactID, created.Etag)
}

func UpdateExistingContact(client *icloud.Client, contactID, etag string) {
	contact := icloud.Contact{
		ContactID: contactID,
		Etag:      etag,
		FirstName: "Jane",
		LastName:  "Doe",
	}

	updated, err := client.UpdateContact(contact)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated contact: %s (New Etag: %s)\n", updated.FirstName, updated.Etag)
}

func DeleteExistingContact(client *icloud.Client, contactID, etag string) {
	err := client.DeleteContact(contactID, etag)
	if err != nil {
		panic(err)
	}
	fmt.Println("Contact deleted successfully")
}
