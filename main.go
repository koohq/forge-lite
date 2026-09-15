package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/koohq/forge-lite/internal/engine"
)

// version is overridden at release build time via -ldflags "-X main.version=...".
var version = ""

func resolveVersion(v string, buildInfoFn func() (*debug.BuildInfo, bool)) string {
	if v != "" {
		return strings.TrimPrefix(v, "v")
	}
	if buildInfoFn != nil {
		if info, ok := buildInfoFn(); ok {
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				return strings.TrimPrefix(info.Main.Version, "v")
			}
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" && setting.Value != "" {
					rev := setting.Value
					if len(rev) > 7 {
						rev = rev[:7]
					}
					return rev
				}
			}
		}
	}
	return "devel"
}

func getVersion() string {
	return resolveVersion(version, debug.ReadBuildInfo)
}

func printHelp() {
	fmt.Println("Forge Lite - A lightweight terminal CLI roguelite")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  forge-lite [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -v, --version  Show version information")
	fmt.Println("  -h, --help     Show help information")
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Println("forge-lite " + getVersion())
			return
		case "-h", "--help":
			printHelp()
			return
		}
	}

	game := engine.NewGame(os.Stdin)
	game.Start()
}
