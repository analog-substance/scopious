package main

import (
	"github.com/analog-substance/scopious/pkg/cmd"
	"github.com/analog-substance/util/cli/build_info"
	"github.com/analog-substance/util/cli/completion"
	"github.com/analog-substance/util/cli/glamour_help"
	"github.com/analog-substance/util/cli/updater/cobra_updater"
)

var version = "v0.0.0"
var commit = "replace"

func main() {
	buildVersion := build_info.GetVersion(version, commit)
	cmd.RootCmd.Version = buildVersion.String()
	cobra_updater.AddToRootCmd(cmd.RootCmd, buildVersion)
	completion.AddToRootCmd(cmd.RootCmd)
	glamour_help.AddToRootCmd(cmd.RootCmd)

	cmd.Execute()
}
