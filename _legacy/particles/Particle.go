package particles

import (
	"errors"
	"fmt"
	"github.com/m41denx/particle-engine/pkg"
	"github.com/m41denx/particle-engine/structs"
	"github.com/m41denx/particle-engine/utils"
	"io"
	"os"
	"path/filepath"
)

var ldir = ""

var homeDir = utils.PrepareStorage()

type Particle struct {
}

func NewParticle() *Particle {
	var err error
	pkg.Config, err = structs.LoadConfig(filepath.Join(homeDir, "config.json"))
	if err != nil {
		pkg.Config.SaveTo(filepath.Join(homeDir, "config.json"))
	}
	return &Particle{}
}

func (p *Particle) GenerateMainfestFor(srvid string, version string) string {
	apply := "2.1"
	base := "2.1"
	if version == "2.2" {
		apply = "2.2"
		base = "2.206"
	}
	return fmt.Sprintf(`
name: %s@v1.0
meta:
    author: ghost
    note: autogen
layer:
    block: '[sha256 autogen]'
recipe:
    - use: fruitspace/gdps_windows@%s
    - apply: fruitspace/msvc_redist@2010
    - apply: fruitspace/msvc_redist@2013
	- apply: fruitspace/gdps_patcher@%s
	  env:
	    GDPS_ID: %s


`, srvid, version, base, apply, srvid)
}

func (p *Particle) InitFolder(path string) {
	particle.ParticleInit(path, "")
}

func (p *Particle) Prepare(path string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	p.p, err = particle.NewParticleFromFile(filepath.Join(path, "particle.json"))
	if err != nil {
		return
	}
	p.p.Analyze(false)
	return
}

func (p *Particle) MoveBuild(path string) (err error) {
	src := filepath.Join(path, ".build", "dist")
	// Move everything from folder src/* to path/
	err = copyRecursively(src, path)
	if err != nil {
		return
	}
	// remove folder src
	os.RemoveAll(filepath.Join(path, ".build"))
	return
}

func copyRecursively(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Calculate the proper destination path for the current item
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			// Create the directory in the dest path
			return os.MkdirAll(destPath, info.Mode())
		} else {
			// If not a directory (i.e., a file), copy it over
			// Open source file
			sourceFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			// Create dest file
			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			// Copy contents to dest file
			if _, err := io.Copy(destFile, sourceFile); err != nil {
				return err
			}
		}
		// After copying is done, set the same permissions as the source
		return os.Chmod(destPath, info.Mode())
	})
}
