package soundcloud

import (
	"fmt"
)

// Host returns the Soundcloud api host address
func Host() string {
	return "https://api.soundcloud.com"
}
// Client ...
func Client(clientID string, clientSecret string) {
	host := Host()
	fmt.Println(host)

}
