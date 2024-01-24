package utils

import (
	"encoding/json"
	"errors"
	"github.com/minio/selfupdate"
	"io"
	"net/http"
	"os"
)

type Server struct {
	Name    string `json:"name"`
	SrvId   string `json:"srvid"`
	Players int    `json:"players"`
	Levels  int    `json:"levels"`
	Icon    string `json:"icon"`
	Version string `json:"version"`
	Recipe  string `json:"recipe"`
}

func LoadServerInfo(srvid string) (Server, error) {
	r, err := http.Get("https://api.fruitspace.one/v2/repatch/gd/" + srvid)
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
