package main

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

var DownloadingNow = make(chan int64)

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {return false}
	return true
}

func CreateWorkdir() string {
	path, _ := os.UserConfigDir()
	os.MkdirAll(path+"/GhostLauncher/configs", 0755)
	os.MkdirAll(path+"/GhostLauncher/cache", 0755)
	return path+"/GhostLauncher"
}

func DownloadDefaultTextures() error {
	pwd:=CreateWorkdir()
	if FileExists(pwd+"/gdps_textures.zip") {
		DownloadingNow<- 1
		return nil
	}
	return DownloadFile("https://cdn.fruitspace.one/assets/gdps_textures.zip",pwd+"/gdps_textures.zip")
}

func DownloadDefaultDLLs() error {
	pwd:=CreateWorkdir()
	if FileExists(pwd+"/gdps_dlls.zip") {
		DownloadingNow <- 1
		return nil
	}
	return DownloadFile("https://cdn.fruitspace.one/assets/gdps_dlls.zip",pwd+"/gdps_dlls.zip")
}




func DownloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {return err}
	defer out.Close()
	//DownloadingNow = make(chan int64)
	resp, err := http.Get(url)
	if err != nil {return err}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	DownloadingNow <- n
	return err
}

// GetWebFileInfo returns filesize and etag
func GetWebFileInfo(url string) (int, string, error) {
	resp, err := http.Head(url)
	if err != nil {return 0, "", err}
	return int(resp.ContentLength),
		strings.ReplaceAll(resp.Header.Get("etag"),"\"",""), nil
}

func GetFileSize(filepath string) (int64, error) {
	file, err := os.Open(filepath)
	if err != nil {return 0, err}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {return 0, err}
	return fi.Size(), nil
}

func GetDownloadPercent(size int, path string) float64 {
	fileSize, _ := GetFileSize(path)
	if fileSize == 0 {fileSize=1}
	return float64(fileSize) / float64(size) * 100
}

func UnzipFile(filepath string, dest string) error {
	r, err := zip.OpenReader(filepath)
	if err != nil {return err}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {return err}
		defer rc.Close()

		path := dest + "/" + f.Name
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {return err}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}


func CheckPath(path string) bool {
	// check if path is ASCII
	for _, r := range path {
		if r > unicode.MaxASCII {return false}
	}
	return true
}


func CheckGDIntegrity() string {
	// get current directory
	pwd, _ := os.Getwd()
	if FileExists(pwd+"/fruit.lock") {return pwd}
	return ""
}


func StartBinaryDetached(path string){
	bin:= exec.Command(path)
	bin.Dir = path[:strings.LastIndex(path, "/")]
	bin.Start()
	bin.Process.Release()
}

func CacheIcon(url string, iconEtag string) error {
	basePath:=CreateWorkdir()
	icon:=url[strings.LastIndex(url, "/")+1:]
	_, etag, _ := GetWebFileInfo(url)
	if etag == iconEtag {return nil}
	go func() {
		for {select {case <-DownloadingNow: break}}
	}()
	err:=DownloadFile(url,basePath+"/cache/"+icon)
	DownloadingNow = make(chan int64)
	return err
}

func CalculateMD5(filepath string) (string, error) {
	f, err := os.Open(filepath)
    if err!= nil {return "", err}
    defer f.Close()
    h := md5.New()
    _, err = io.Copy(h, f)
	if err!= nil {return "", err}
	return fmt.Sprintf("%x",h.Sum(nil)), nil
}