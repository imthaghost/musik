package mp3

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/grafov/m3u8"
	"github.com/imthaghost/musik/decrypt"
	"github.com/imthaghost/musik/joiner"
	"github.com/imthaghost/musik/pool"
	"github.com/imthaghost/musik/zhttp"
)

var (
	// ZHTTP client
	ZHTTP *zhttp.Zhttp
	// JOINER client
	JOINER       *joiner.Joiner
	keyCache     map[string][]byte
	keyCacheLock sync.Mutex
)

func start(mpl *m3u8.MediaPlaylist) {
	p := pool.New(10, download)

	go func() {
		var count = int(mpl.Count())
		for i := 0; i < count; i++ {
			p.Push([]interface{}{i, mpl.Segments[i], mpl.Key})
		}
		p.CloseQueue()
	}()

	go p.Run()
}

func parseM3u8(m3u8Url string) (*m3u8.MediaPlaylist, error) {
	statusCode, data, err := ZHTTP.Get(m3u8Url)
	if err != nil {
		return nil, err
	}

	if statusCode/100 != 2 || len(data) == 0 {
		return nil, errors.New("download m3u8 file failed, http code: " + strconv.Itoa(statusCode))
	}

	playlist, listType, err := m3u8.Decode(*bytes.NewBuffer(data), true)
	if err != nil {
		return nil, err
	}

	if listType == m3u8.MEDIA {
		obj, _ := url.Parse(m3u8Url)
		mpl := playlist.(*m3u8.MediaPlaylist)

		if mpl.Key != nil && mpl.Key.URI != "" {
			uri, err := formatURI(obj, mpl.Key.URI)
			if err != nil {
				return nil, err
			}
			mpl.Key.URI = uri
		}

		count := int(mpl.Count())
		for i := 0; i < count; i++ {
			segment := mpl.Segments[i]

			uri, err := formatURI(obj, segment.URI)
			if err != nil {
				return nil, err
			}
			segment.URI = uri

			if segment.Key != nil && segment.Key.URI != "" {
				uri, err := formatURI(obj, segment.Key.URI)
				if err != nil {
					return nil, err
				}
				segment.Key.URI = uri
			}

			mpl.Segments[i] = segment
		}

		return mpl, nil
	}

	return nil, errors.New("Unsupport m3u8 type")
}

func getKey(url string) ([]byte, error) {
	keyCacheLock.Lock()
	defer keyCacheLock.Unlock()

	key := keyCache[url]
	if key != nil {
		return key, nil
	}

	statusCode, key, err := ZHTTP.Get(url)
	if err != nil {
		return nil, err
	}

	if len(key) == 0 {
		return nil, errors.New("body is empty, http code: " + strconv.Itoa(statusCode))
	}

	keyCache[url] = key

	return key, nil
}

func download(in interface{}) {
	params := in.([]interface{})
	id := params[0].(int)
	segment := params[1].(*m3u8.MediaSegment)
	globalKey := params[2].(*m3u8.Key)

	statusCode, data, err := ZHTTP.Get(segment.URI)
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s Download failed: %s\n", red("[-]"), err)
	}

	if len(data) == 0 {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s Download failed: body is empty, http code: %d\n", red("[-]"), statusCode)
	}

	var keyURL, ivStr string
	if segment.Key != nil && segment.Key.URI != "" {
		keyURL = segment.Key.URI
		ivStr = segment.Key.IV
	} else if globalKey != nil && globalKey.URI != "" {
		keyURL = globalKey.URI
		ivStr = globalKey.IV
	}

	if keyURL != "" {
		var key, iv []byte
		key, err = getKey(keyURL)
		if err != nil {
			fmt.Println("[-] Download key failed:", keyURL, err)
		}

		if ivStr != "" {
			iv, err = hex.DecodeString(strings.TrimPrefix(ivStr, "0x"))
			if err != nil {
				fmt.Println("[-] Decode iv failed:", err)
			}
		} else {
			iv = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(id)}
		}

		data, err = decrypt.Decrypt(data, key, iv)
		if err != nil {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s Decrypt failed: %s\n", red("[-]"), err)
		}
	}

	// log.Println("[+] Download succed:", segment.URI)

	JOINER.Join(id, data)
}

func formatURI(base *url.URL, u string) (string, error) {
	if strings.HasPrefix(u, "http") {
		return u, nil
	}

	obj, err := base.Parse(u)
	if err != nil {
		return "", err
	}

	return obj.String(), nil
}

func filename(u string) string {
	obj, _ := url.Parse(u)
	_, filename := filepath.Split(obj.Path)
	return filename
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// Merge ...
func Merge(url string) string {

	keyCache = map[string][]byte{}

	var err error
	ZHTTP, err = zhttp.New(time.Second*30, "")
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s Init failed: %s\n", red("[-]"), err)
	}

	mpl, err := parseM3u8(url)
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s Parse m3u8 file failed: %s\n", red("[-]"), err)
	} else {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Parse m3u8 file succed %s\n", green("[+]"), "")
	}
	name := RandStringBytesMaskImpr(5)

	outFile := name + ".mp3"

	JOINER, err = joiner.New(outFile)

	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s Open file failed: %s\n", red("[-]"), err)
	} else {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Will save to %s\n", green("[+]"), JOINER.Name())
	}

	if mpl.Count() > 0 {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Total %d files to download \n", green("[+]"), mpl.Count())

		start(mpl)

		err = JOINER.Run(int(mpl.Count()))
		if err != nil {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s Write to file failed: %s\n", red("[-]"), err)
		}
		g := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Download succed, saved to %s\n", g("[+]"), JOINER.Name())

	}
	return outFile
}
