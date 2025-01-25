module GhostPatcher

go 1.23

require (
	fyne.io/fyne/v2 v2.4.3
	github.com/fyne-io/terminal v0.0.0-20240206170039-2e129cdfd85f
	github.com/klauspost/cpuid/v2 v2.2.6
	github.com/m41denx/particle-engine v0.10.3-beta
	github.com/minio/selfupdate v0.6.0
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58
	golang.org/x/sys v0.27.0
)

require (
	aead.dev/minisign v0.3.0 // indirect
	fyne.io/systray v1.10.1-0.20231115130155-104f5ef7839e // indirect
	github.com/ActiveState/termtest/conpty v0.5.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/creack/pty v1.1.21 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fredbi/uri v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fyne-io/gl-js v0.0.0-20230506162202-1fdaa286a934 // indirect
	github.com/fyne-io/glfw-js v0.0.0-20240101223322-6e1efdc71b7a // indirect
	github.com/fyne-io/image v0.0.0-20240121103648-c3c798e60e6b // indirect
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240118000515-a250818d05e3 // indirect
	github.com/go-text/render v0.0.0-20240122202426-67aad72d5803 // indirect
	github.com/go-text/typesetting v0.1.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jsummers/gobmp v0.0.0-20230614200233-a9de23ed2e25 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/tevino/abool v1.2.0 // indirect
	github.com/yuin/goldmark v1.7.0 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/image v0.15.0 // indirect
	golang.org/x/mobile v0.0.0-20240112133503-c713f31d574b // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/js/dom v0.0.0-20231112215516-51f43a291193 // indirect
)

replace (
	github.com/go-text/render => github.com/go-text/render v0.0.0-20230619120952-35bccb6164b8
	github.com/go-text/typesetting => github.com/go-text/typesetting v0.0.0-20230616162802-9c17dd34aa4a
	github.com/m41denx/particle => github.com/m41denx/particle-engine v0.10.3-beta
)
