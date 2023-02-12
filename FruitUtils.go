package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minio/selfupdate"
	"io"
	"net/http"
	"os"
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

func LoadServerInfo(srvid string) (Server, error) {
	r, err := http.Get("https://api.fruitspace.one/v1/repatch/getserverinfo?id=" + srvid)
	if err != nil {
		return Server{}, err
	}
	if r.StatusCode != 200 {
		return Server{}, errors.New("[" + srvid + "] Сервер не найден")
	}
	defer r.Body.Close()
	var s Server
	err = json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		return Server{}, err
	}
	return s, nil
}

func (s Server) GetUrl() string {
	return "http://" + s.Region + "gd.gofruit.space/" + s.SrvId + "/db"
}

type RePatcher struct{}

func (rp RePatcher) DownloadPureGD() ([]byte, error) {
	resp, err := http.Get("https://cdn.fruitspace.one/assets/GeometryDash.exe")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// PatchPureGD url is http://XXX/db/
func (rp RePatcher) PatchPureGD(url string, gd []byte) []byte {
	gd = bytes.ReplaceAll(gd, []byte("http://www.boomlings.com/database"), []byte(url))
	gd = bytes.ReplaceAll(gd, []byte("RobTop Support for more info"), []byte("Fruitspace Support for help."))
	gd = bytes.ReplaceAll(gd, []byte("Something went wrong\nplease try again later"), []byte("Nothing here yet :/ \nmaybe try again later?"))
	encoded := base64.StdEncoding.EncodeToString([]byte(url))
	encoded = MinifyBase64(encoded)
	gd = bytes.ReplaceAll(gd, []byte("aHR0cDovL3d3dy5ib29tbGluZ3MuY29tL2RhdGFiYXNl"), []byte(encoded))
	return gd
}

func MinifyBase64(data string) string {
	if len(data) > 46 {
		if strings.HasSuffix(data, "w==") {
			return data[:45] + "3"
		} else {
			fmt.Printf("[Minify] Unable to patch (%d/46)\n", len(data))
		}
	}
	return data
}

func WriteBytes(filepath string, data []byte) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func Update(path string) error {
	r, err := http.Get("https://cdn.fruitspace.one/assets/GhostLauncher.exe")
	if err != nil {
		return err
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return WriteBytes(path, data)
}

func SelfUpdate() error {
	r, err := http.Get("https://cdn.fruitspace.one/assets/GhostLauncher.exe")
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return selfupdate.Apply(r.Body, selfupdate.Options{})
}

type Lock struct {
	SrvId              string `json:"srvid"`
	Title              string `json:"title"`
	Uname              string `json:"uname"`
	TextureEtag        string `json:"textureEtag"`
	TextureOverlayEtag string `json:"textureOverlayEtag"`
	DllEtag            string `json:"dllEtag"`
	IconEtag           string `json:"iconEtag"`
	LauncherEtag       string `json:"launcherEtag"`

	ModLoaderVersion string   `json:"modLoaderVersion"`
	Mods             []string `json:"mods"`
}

func (l *Lock) WriteLock(path string) error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return WriteBytes(path+"/fruit.lock", data)
}

func (l *Lock) ReadLock(path string) error {
	f, err := os.Open(path + "/fruit.lock")
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&l)
	return err
}

type ModUnit struct {
	Id          string `json:"id"`
	Target      string `json:"target"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	HasTextures bool   `json:"hasTextures"`
	HasConfig   bool   `json:"hasConfig"`
}
