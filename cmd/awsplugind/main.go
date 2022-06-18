package main

import (
	"strings"

	"go.autokitteh.dev/sdk/pluginsvc"

	"github.com/autokitteh/aws/internal/pkg/plugin"
)

var version, commit, date string

func init() {
	version = strings.TrimPrefix(version, "v")
}

func main() {
	pluginsvc.Run(
		&pluginsvc.Version{Version: version, Commit: commit, Date: date},
		plugin.NewPlugin(),
	)
}
