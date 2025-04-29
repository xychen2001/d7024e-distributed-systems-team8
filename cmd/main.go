package main

import (
	"github.com/eislab-cps/go-template/internal/cli"
	"github.com/eislab-cps/go-template/pkg/build"
)

var (
	BuildVersion string = ""
	BuildTime    string = ""
)

func main() {
	build.BuildVersion = BuildVersion
	build.BuildTime = BuildTime
	cli.Execute()
}
