package main

import (
	"github.com/qlik-oss/corectl/cmd"
)

// version will be set with: go build -ldflags "-X main.version=X.Y.Z"
var version = "development build"

func main() {
	cmd.Execute(version)
}
