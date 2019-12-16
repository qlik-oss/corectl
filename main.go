package main

import (
	"github.com/qlik-oss/corectl/cmd"
)

// version will be set with: go build -ldflags "-X main.version=X.Y.Z"
var version = ""
var commit = ""
var branch = ""

func main() {
	cmd.Execute(version, branch, commit)
}
