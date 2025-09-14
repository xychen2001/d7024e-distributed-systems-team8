package main

import (
	"github.com/xychen/d7024e-distributed-systems-team8/internal/cli"
	"github.com/xychen/d7024e-distributed-systems-team8/pkg/build"
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
