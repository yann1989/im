// Author       kevin
// Time         2019-08-08 20:16
// File Desc    版本管理

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli"
)

const (
	VersionMajor = 1          // 重大版本
	VersionMinor = 0          // 次要版本
	VersionPatch = 0          // 版本补丁
	VersionMeta  = "unstable" // 是否稳定
)

// 格式化版本信息
var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

// 版本元数据信息
var VersionWithMeta = func() string {
	v := Version
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

func version(_ *cli.Context) string {
	ver := fmt.Sprintln(strings.Title(clientIdentifier))

	ver += fmt.Sprintln("Version:", VersionWithMeta)
	if gitCommit != "" {
		ver += fmt.Sprintln("Git Commit:", gitCommit)
	}
	ver += fmt.Sprintln("Architecture:", runtime.GOARCH)
	ver += fmt.Sprintln("Go Version:", runtime.Version())
	ver += fmt.Sprintln("Operating System:", runtime.GOOS)
	ver += fmt.Sprintln("GOPATH=", os.Getenv("GOPATH"))
	ver += fmt.Sprintf("GOROOT=%s", runtime.GOROOT())
	return ver
}
