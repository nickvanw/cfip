package cfip

import (
	"encoding/json"
	"net/http"
)

// JSONIP is the URL for where we get our v4 address
const JSONIP = "https://ipv4.jsonip.com/"

// FetchIP is a very tightly scoped and hardcoded method to get the IPv4
// address of your default route to the internet.
func FetchIP() (string, error) {
	resp, err := http.Get(JSONIP)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	var data IP
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.IP, nil
}

// IP has an IP
type IP struct {
	IP string `json:"ip"`
}
