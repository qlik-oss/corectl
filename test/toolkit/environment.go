package toolkit

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var update = true // *flag.Bool("update", false, "update golden files")

var Engine0IP = flag.String("engineStd", "localhost:9076", "URL to first engine instance in docker-compose.yml i.e qix-engine-std")
var Engine1IP = flag.String("engineJwt", "localhost:9176", "URL to first engine instance in docker-compose.yml i.e qix-engine-jwt")
var Engine2IP = flag.String("engineAbac", "localhost:9276", "URL to third engine instance in docker-compose.yml i.e qix-engine-abac")
var Engine3IP = flag.String("engineBadLicenseServer", "localhost:9376", "URL to second engine instance in docker-compose.yml i.e qix-engine-bad-license-server")

func init() {
	fmt.Println("RUNNING TEST MAIN")
	buildCorectl()

	os.Setenv("CORECTL_TEST_CONNECT", "corectl-test-connector")
	os.Setenv("ENGINE_STD_URL", *Engine0IP)

	AddGoldPolisher("(New connection created with id): .*$", "$1: <filtered for gold shininess>")
}

func buildCorectl() {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}

	abs, err := filepath.Abs(binaryName)
	if err != nil {
		fmt.Printf("could not get abs path for %s: %v", binaryName, err)
		os.Exit(1)
	}

	binaryPath = abs

	if err := exec.Command("go", "build", "-o", binaryName, "-v").Run(); err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}
}

func getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "corectl.exe"
	}
	return "corectl"
}

var binaryName = getBinaryName()

var binaryPath string

func GetTestFilePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
