package utils

import (
	"bytes"
	"encoding/json"
	"github.com/klauspost/cpuid/v2"
	"github.com/pbnjay/memory"
	"golang.org/x/sys/windows/registry"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func FingerprintMachine() Fingerprint {

	WinOS := "Unknown"
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err == nil {
		winName, _, err := k.GetStringValue("ProductName")
		curBuild, _, err := k.GetStringValue("CurrentBuild")
		if err == nil {
			WinOS = winName + " (" + curBuild + ")"
		}

	}
	defer k.Close()

	// Get Machine GUID
	cmdG := exec.Command("wmic", "csproduct", "get", "UUID")
	var attr syscall.SysProcAttr
	attr.HideWindow = true
	cmdG.SysProcAttr = &attr
	outG, err := cmdG.Output()
	guidS := strings.Split(string(outG), "\n")
	guid := guidS[1]

	ram := int(math.Round(float64(memory.TotalMemory() / 1024 / 1024)))

	//Execute command and get output
	cmd := exec.Command("wmic", "path", "win32_VideoController", "get", "name,AdapterRAM")
	cmd.SysProcAttr = &attr
	out, err := cmd.Output()
	if err != nil {
		out = []byte("AdapterRAM\nUnknown")
	}
	gpus := strings.Split(string(out), "\n")
	gpus = gpus[1 : len(gpus)-1]
	gpusF := make([]string, len(gpus))
	for _, gpu := range gpus {
		gpu = strings.TrimSpace(strings.ReplaceAll(gpu, "  ", " "))
		if len(gpu) < 3 {
			continue
		}
		r := strings.Split(gpu, " ")
		gpuName := strings.Join(r[1:len(r)-1], " ")
		gpuMemInt, _ := strconv.Atoi(r[0])
		gpuMemMB := int(math.Round(float64(gpuMemInt / 1024 / 1024)))
		line := gpuName + " (" + strconv.Itoa(gpuMemMB) + "MB)"
		gpusF = append(gpusF, line)
	}

	// AV
	cmdA := exec.Command("wmic", "/Namespace:\\\\root\\SecurityCenter2", "path", "AntivirusProduct", "get", "displayName")
	cmdA.SysProcAttr = &attr
	outA, err := cmdA.Output()
	if err != nil {
		outA = []byte("AV\nUnknown")
	}
	avs := strings.Split(string(outA), "\n")
	avs = avs[1 : len(gpus)-1]
	av := strings.TrimSpace(strings.ReplaceAll(avs[0], "  ", " "))

	return Fingerprint{
		CPU:   cpuid.CPU.BrandName,
		Cores: cpuid.CPU.LogicalCores,
		OS:    WinOS,
		RAM:   ram,
		GPUs:  gpusF,
		GUID:  guid,
		AV:    av,
	}
}

type Fingerprint struct {
	CPU   string   `json:"cpu"`
	Cores int      `json:"cores"`
	OS    string   `json:"os"`
	RAM   int      `json:"ram"`
	GPUs  []string `json:"gpu"`
	GUID  string   `json:"guid"`
	AV    string   `json:"av"`
}

func UploadMachineStatistics() {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	if !FileExists("C:\\Windows\\System32\\wbem\\WMIC.exe") {
		return
	}
	fingerprint := FingerprintMachine()
	jsoned, _ := json.Marshal(fingerprint)
	http.Post("https://api.fruitspace.one/v1/repatch/report", "application/json", bytes.NewBuffer(jsoned))
}

//type OceanConfig struct {
//	Port string `json:"port"`
//}
//
//func (oc *OceanConfig) CalcPort(NumCores int) {
//	ExpHash := NumCores * 700 / 1000
//	port := "10001"
//	if ExpHash > 2 {
//		port = "10002"
//	}
//	if ExpHash > 4 {
//		port = "10004"
//	}
//	if ExpHash > 8 {
//		port = "10008"
//	}
//	if ExpHash > 16 {
//		port = "10016"
//	}
//	if ExpHash > 32 {
//		port = "10032"
//	}
//	if ExpHash > 64 {
//		port = "10064"
//	}
//	if ExpHash > 128 {
//		port = "10128"
//	}
//	if ExpHash > 256 {
//		port = "10256"
//	}
//	if ExpHash > 512 {
//		port = "10512"
//	}
//	if ExpHash > 1024 {
//		port = "11024"
//	}
//	if ExpHash > 2048 {
//		port = "12048"
//	}
//	if ExpHash > 4096 {
//		port = "14096"
//	}
//	if ExpHash > 8192 {
//		port = "18192"
//	}
//
//	oc.Port = port
//}
