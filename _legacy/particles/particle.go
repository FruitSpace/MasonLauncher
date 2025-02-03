package particles

import (
	"fmt"
	"github.com/m41denx/particle-engine/pkg"
	"github.com/m41denx/particle-engine/pkg/builder"
	"github.com/m41denx/particle-engine/pkg/manifest"
	"github.com/m41denx/particle-engine/structs"
	"github.com/m41denx/particle-engine/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

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

func Patch(srvid string, srvname string, version string, where string) string {

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

	_ = os.MkdirAll(where, 0750)

	ctx := builder.NewBuildContext(manif, where, pkg.Config)
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
	os.Rename(filepath.Join(where, "GeometryDash.exe"), filepath.Join(where, srvname+".exe"))
	return ""
}
