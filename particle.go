package main

import (
	"fmt"
	"github.com/m41denx/particle-engine/pkg"
	"github.com/m41denx/particle-engine/pkg/builder"
	"github.com/m41denx/particle-engine/pkg/manifest"
	"github.com/m41denx/particle-engine/structs"
	"github.com/m41denx/particle-engine/utils"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var piper, pipew, _ = os.Pipe()

func GetRecipe(version string, patcher string, srvid string) string {
	return fmt.Sprintf(`
name: %s_gdps
meta:
  author: mason
  note: 'Its a yaml, dont you fucking touch it, idiot'
layer:
  block: '[sha256 autogen]'
recipe:
  - use: fruitspace/gdps_windows@%s
  - apply: fruitspace/msvc_redist@2010
  - apply: fruitspace/msvc_redist@2013
  - apply: fruitspace/gdps_patcher
    env:
      GDPS_ID: %s
      GDPS_VER: %s
`, srvid, version, srvid, patcher)
}

var homeDir = utils.PrepareStorage()
var pc = filepath.Join(fsRootDir(), "MasonLauncher")

func PrepareLauncher() string {
	_ = os.MkdirAll(pc, 0750)
	return pc
}

func (*App) Read() string {
	buf := make([]byte, 1024)
	piper.Read(buf)
	return string(buf)
}

func (*App) Patch(srvid string, srvname string, version string) string {
	os.Stdout = pipew
	os.Stderr = pipew

	ver := "2.204"
	if version != "2.2" {
		ver = version
	}
	recipe := GetRecipe(ver, version, srvid)

	// Start particle ss
	var err error
	pkg.Config, err = structs.LoadConfig(filepath.Join(homeDir, "config.json"))
	if err != nil {
		pkg.Config.SaveTo(filepath.Join(homeDir, "config.json"))
	}

	var manif manifest.Manifest
	err = yaml.Unmarshal([]byte(recipe), &manif)
	if err != nil {
		return err.Error()
	}

	_ = os.MkdirAll(filepath.Join(pc, srvid), 0750)

	ctx := builder.NewBuildContext(manif, filepath.Join(pc, srvid), pkg.Config)
	fmt.Println("Starting build")
	if err := ctx.FetchDependencies(); err != nil {
		ctx.Clean(false)
		return err.Error()
	}
	if err := ctx.PrepareEnvironment(); err != nil {
		ctx.Clean(false)
		return err.Error()
	}
	if err = ctx.Export(); err != nil {
		return err.Error()
	}
	ctx.Clean(false)
	os.Rename(filepath.Join(pc, srvid, "GeometryDash.exe"), filepath.Join(pc, srvid, srvname+".exe"))
	return ""
}

func (*App) StartGDPS(srvid string, srvname string) {
	path := filepath.Join(pc, srvid, srvname+".exe")
	bin := exec.Command(path)
	bin.Dir = path[:strings.LastIndex(path, "\\")]
	if err := bin.Start(); err != nil {
		fmt.Println(err)
		return
	}
	bin.Process.Release()
}

func fsRootDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("SystemDrive") + "\\"
	}
	hen, err := os.UserHomeDir()
	if err != nil {
		hen = "/"
	}
	return hen
}

func (*App) ListServers() (m []string) {
	l, e := os.ReadDir(pc)
	if e != nil {
		return
	}
	for _, f := range l {
		if f.IsDir() {
			m = append(m, f.Name())
		}
	}
	return
}
