package soundcloud

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

func RandomSong() {
	rand.Seed(time.Now().Unix())
	// top playlists
	hits := [...]string{"https://soundcloud.com/soundcloud-the-peak/sets/on-the-up-the-peak-hot-new", "https://soundcloud.com/soundcloud-shine/sets/new-now", "https://soundcloud.com/soundcloud-hustle/sets/hip-hop-party-starters", "https://soundcloud.com/soundcloud-the-peak/sets/chill-edm"}
	// random integer
	n := rand.Int() % len(hits)
	// request random playlist
	resp, err := http.Get(hits[n])
	if err != nil {
		log.Fatalln(err)
	}
	// response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// // if the song that we found from response body is not in the songs list add it otherwise continue
	// found := contains(songs, "B")
	// if !found {
	// 	fmt.Println("Value not found in slice")
	// }

	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	list := htmlquery.Find(doc, "/html/body/script[9]/text()")
	var songlist []string
	for _, n := range list {
		fmt.Println(htmlquery.InnerText(n))
		songlist = append(songlist, htmlquery.InnerText(n))
	}

}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
