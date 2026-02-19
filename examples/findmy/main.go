package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Johnw7789/Go-iClient/icloud"
)

func main() {
	// client := newClient("username", "password")

	// Uncomment the example you want to run
	// GetAllDevices(client)
	// GetSingleDevice(client, "device-id")
	// PlaySoundOnDevice(client, "device-id", nil)
	// PlaySoundOnDevice(client, "airpods-device-id", []string{"left", "right"})
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

// GetAllDevices prints all devices on the account including family members
func GetAllDevices(client *icloud.Client) {
	devices, err := client.GetDevices()
	if err != nil {
		panic(err)
	}

	for _, d := range devices {
		fmt.Printf("[%s] %s — %s (status: %s)\n\tid: %s\n", d.DeviceClass, d.Name, d.DeviceDisplayName, d.DeviceStatus, d.ID)
	}
}

// GetSingleDevice prints details for a specific device
func GetSingleDevice(client *icloud.Client, deviceID string) {
	device, err := client.GetDevice(deviceID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name:         %s\n", device.Name)
	fmt.Printf("Display Name: %s\n", device.DeviceDisplayName)
	fmt.Printf("Class:        %s\n", device.DeviceClass)
	fmt.Printf("Model:        %s\n", device.RawDeviceModel)
	fmt.Printf("Status:       %s\n", device.DeviceStatus)
	fmt.Printf("Lost Mode:    %v\n", device.LostModeEnabled)
	if device.BatteryStatus != nil {
		fmt.Printf("Battery:      %s\n", *device.BatteryStatus)
	}
}

// PlaySoundOnDevice triggers a sound alert on a device
// pass nil channels for iPhone/Watch/Mac, []string{"left", "right"} for AirPods
func PlaySoundOnDevice(client *icloud.Client, deviceID string, channels []string) {
	// keepalive in case the operation takes a while
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := client.KeepAlive(ctx, 5*time.Minute); err != nil && !errors.Is(err, context.Canceled) {
			fmt.Println("Session keepalive stopped:", err)
		}
	}()

	updated, err := client.PlaySound(deviceID, channels)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sound triggered on %s (snd status: %s)\n", updated.Name, updated.Snd.StatusCode)
}
