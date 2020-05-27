// Author       kevin
// Time         2019-08-08 20:16
// File Desc    主入口

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	clientIdentifier = "Taurus Chicha Service"
)

var (
	// Git SHA1提交发布的哈希值
	gitCommit = ""
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	startCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Please use ./"+filepath.Base(os.Args[0])+" start --config config.toml")
	rootCmd.AddCommand(startCmd)
}

func main() {
	version(nil)
	// 执行主逻辑
	if err := Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
