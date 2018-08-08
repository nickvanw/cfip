package main

import (
	"fmt"
	"os"

	cfip "github.com/nickvanw/cfip"
)

func main() {
	client, err := cfip.NewClient(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"), os.Getenv("CF_ZONE"))
	if err != nil {
		fmt.Printf("failed to create client: %s\n", err)
		os.Exit(1)
	}
	ip, err := cfip.FetchIP()
	if err != nil {
		fmt.Printf("failed to get IP: %s\n", err)
		os.Exit(2)
	}
	if err := client.Set(os.Getenv("CF_HOST"), ip); err != nil {
		fmt.Printf("failed to update IP: %s\n", err)
		os.Exit(3)
	}
}
