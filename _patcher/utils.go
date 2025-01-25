package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
)

type Server struct {
	Name        string `json:"name"`
	SrvId       string `json:"srvid"`
	Players     int    `json:"players"`
	Levels      int    `json:"levels"`
	Icon        string `json:"icon"`
	TexturePack string `json:"texturepack"`
	Region      string `json:"region"`
}

func GetUrl(srvid string) string {
	return fmt.Sprintf("http://rugd.gofruit.space/%s/db", srvid)
}

var ImmutableURL = func(rot int, s bool) (string, func() string, []byte) {
	pre := fmt.Sprintf("http%s://www.", t(s, "s", ""))
	data := struct {
		Sound string
		Name  []byte
		Data  string
		Base  []byte
	}{"boom", []byte("lings"), "data", []byte("base")}
	return pre + data.Sound, func() string { return string(data.Name) + ".com/" + data.Data }, data.Base
}

func PatchPureGD(url string, gd []byte, s bool) []byte {

	legacyUrl := url
	if s {
		url = strings.ReplaceAll(url, "http://", "https://")
	}

	fmt.Printf("Patching GPDS URL: %s\n", url)

	a, b, c := ImmutableURL(15, s)
	oldUrl := []byte(a + b() + string(c))
	newUrl := []byte(url)

	fmt.Printf("Will replace:\n%s\n%s\n", string(oldUrl), string(newUrl))

	gd = bytes.ReplaceAll(gd, oldUrl, newUrl)

	oldEncodedUrl := []byte("aHR0cDovL3d3dy5ib29tbGluZ3MuY29tL2RhdGFiYXNl")
	//if s {
	//	oldEncodedUrl = []byte("aHR0cHM6Ly93d3cuYm9vbWxpbmdzLmNvbS9kYXRhYmFzZS8=")
	//}
	encoded := base64.StdEncoding.EncodeToString([]byte(legacyUrl))
	encoded = minifyBase64(encoded)
	newEncodedUrl := []byte(encoded)

	gd = bytes.ReplaceAll(gd, oldEncodedUrl, newEncodedUrl)

	gd = bytes.ReplaceAll(gd, []byte("RobTop Support for more info"), []byte("FruitSpace Support for help."))
	gd = bytes.ReplaceAll(gd, []byte("Something went wrong\nplease try again later"), []byte("Nothing here yet :/ \nmaybe try again later?"))

	return gd
}
func minifyBase64(data string) string {
	if len(data) > 46 {
		if strings.HasSuffix(data, "w==") {
			return data[:45] + "3"
		}
	}
	return data
}

func t[T any](s bool, a, b T) T {
	if s {
		return a
	} else {
		return b
	}
}
