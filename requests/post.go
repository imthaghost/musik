package requests

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetToken ...
func GetToken() {

	url := "http://3.18.215.82:8000/create"
	method := "POST"

	payload := strings.NewReader("")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	resourceURL := "https://accounts.spotify.com/api/token"
	ClientSecret := "xxxx"
	ClientID := "xxxx"
	clientSJSON := b64.StdEncoding.EncodeToString([]byte(ClientSecret))
	clitentIJSON := b64.StdEncoding.EncodeToString([]byte(ClientID))
	parsed = "Basic" + clientSJSON + clitentIJSON
	credentials := map[string]string{"Authorization": parsed}
	fmt.Println(credentials)
	fmt.Println(clientSJSON)
}
