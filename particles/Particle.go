package particles

import (
	"errors"
	"fmt"
	"github.com/m41denx/particle/particle"
	"github.com/m41denx/particle/utils"
	"github.com/m41denx/particle/utils/hget"
	"io"
	"os"
	"path/filepath"
)

var ldir = ""

type Particle struct {
	p *particle.Particle
}

func NewParticle() *Particle {
	particle.ParticleCache = make(map[string]*particle.Particle)
	particle.LayerCache = make(map[string]*particle.Layer)
	particle.EngineCache = make(map[string]*particle.Engine)
	particle.MetaCache = make(map[string]string)
	particle.NUMCPU = 1

	hget.DisplayProgress = false

	ldir = utils.PrepareStorage()
	return &Particle{}
}

func (p *Particle) GenerateMainfestFor(srvid string, version string) string {
	apply := "2.1"
	base := "2.1"
	if version == "2.2" {
		apply = "2.2"
		base = "2.204"
	}
	return fmt.Sprintf(`
{
        "name": "%s@1.0",
        "author": "ghost",
        "note": "",
        "block": "",
        "meta": {
                "FS_GDPS": "%s",
                "FS_GDPS_VER": "%s"
        },
        "recipe": {
                "base": "m41den/gdps_windows@%s",
                "apply": ["m41den/gdps_patcher@%s"],
                "engines": [],
                "run": []
        }
}
`, srvid, srvid, version, base, apply)
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
	p.p.Analyze()
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
