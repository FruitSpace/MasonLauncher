package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
)

type SaveManager struct {
	Savepath string
	Data     []byte
}

func (sm *SaveManager) Open(srvname string) error {
	local, _ := os.UserCacheDir()
	sm.Savepath = local + "/" + srvname + "/CCGameManager.dat"
	if !FileExists(sm.Savepath) {
		return errors.New("No savedata")
	}
	sd, err := os.Open(sm.Savepath)
	defer sd.Close()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(sd)
	data = XOR(data)
	data = bytes.Trim(data, "\x00")
	vdata := string(data)
	data2, err := base64.URLEncoding.DecodeString(vdata)
	if err != nil {
		return err
	}
	data, err = GzipDecompress(data2)
	if err != nil {
		return err
	}
	sm.Data = data
	return err
}

func (sm *SaveManager) GetUname() string {
	r := bytes.Split(sm.Data, []byte("playerName</k><s>"))[1]
	r = bytes.Split(r, []byte("</s>"))[0]
	return string(r)
}

func XOR(data []byte) []byte {
	r := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		r[i] = data[i] ^ 11
	}
	return r
}

func GzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}
