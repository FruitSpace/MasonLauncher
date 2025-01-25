package main

import (
	"fmt"
	"os"
	"strings"
)

var GDPS_ID = os.Getenv("GDPS_ID")
var GDPS_URL = os.Getenv("GDPS_URL")
var GDPS_VER = os.Getenv("GDPS_VER")

func init() {
	if len(GDPS_ID) == 0 && len(GDPS_URL) == 0 {
		fmt.Println("Error: both GDPS_ID and GDPS_URL environment variables are not set.")
		os.Exit(1)
	}

	if len(GDPS_ID) != 4 {
		if !strings.Contains(GDPS_URL, "http") {
			fmt.Println("Error: GDPS_ID must be exactly 4 characters long.")
			os.Exit(1)
		}
	}

	if len(GDPS_VER) == 0 {
		fmt.Println("Warning: GDPS_VER environment variable is not set, assuming 2.2 (possible: 2.1/2.2)")
		GDPS_VER = "2.2"
	}

	if len(GDPS_ID) == 4 {
		GDPS_URL = GetUrl(GDPS_ID)
	}

}

const VER = "1.0.1"

func main() {

	fmt.Printf("Running GDPS Patcher v%s\n", VER)

	path := fmt.Sprintf("%s/GeometryDash.exe", os.Getenv("BUILD"))

	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: cannot read file: %v\n", err)
		os.Exit(1)
	}

	nbytes := PatchPureGD(GDPS_URL, bytes, GDPS_VER == "2.2")

	err = os.WriteFile(path, nbytes, 0777)
	if err != nil {
		fmt.Printf("Error: cannot write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("GDPS Patcher Completed Successfully")
}
