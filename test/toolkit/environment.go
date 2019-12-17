package toolkit

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var binaryPath string
var update = flag.Bool("update", false, "update golden files")
var EngineStdIP = flag.String("engineStd", "localhost:9076", "URL to first engine instance in docker-compose.yml i.e qix-engine-std")
var EngineJwtIP = flag.String("engineJwt", "localhost:9176", "URL to first engine instance in docker-compose.yml i.e qix-engine-jwt")
var EngineAbacIP = flag.String("engineAbac", "localhost:9276", "URL to third engine instance in docker-compose.yml i.e qix-engine-abac")
var EngineBadLicenseServerIP = flag.String("engineBadLicenseServer", "localhost:9376", "URL to second engine instance in docker-compose.yml i.e qix-engine-bad-license-server")

func init() {
	buildCorectl()
	os.Setenv("CORECTL_TEST_CONNECT", "corectl-test-connector")
	os.Setenv("ENGINE_STD_URL", *EngineStdIP)
	os.Setenv("ENGINE_JWT_URL", *EngineJwtIP)
	os.Setenv("ENGINE_ABAC_URL", *EngineAbacIP)
	os.Setenv("ENGINE_BAD_LICENSE_SERVER_URL", *EngineBadLicenseServerIP)
	AddGoldPolisher("(New connection created with id): .*$", "$1: <filtered for gold shininess>")
	AddGoldPolisher("localhost:9076", "<host>:<port>")
	AddGoldPolisher("(\"qUtcModifyTime\":) .*$", "$1 <filtered for gold shininess>")
}

func buildCorectl() {
	var binaryName string
	if runtime.GOOS == "windows" {
		binaryName = "corectl.exe"
	} else {
		binaryName = "corectl"
	}
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

	args := []string{"build", "-ldflags", "-X main.version=dev", "-o", binaryName, "-v"}
	if err := exec.Command("go", args...).Run(); err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}
}
