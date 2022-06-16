package plugin

import (
	"context"

	"go.autokitteh.dev/sdk/api/apivalues"
	"go.autokitteh.dev/sdk/pluginimpl"
)

var Plugin = &pluginimpl.Plugin{
	ID:  "aws",
	Doc: "AWS SDK Plugin",
	Members: map[string]*pluginimpl.PluginMember{
		"cat": pluginimpl.NewSimpleMethodMember(
			"returns cat's vocalization",
			func(context.Context, []*apivalues.Value, map[string]*apivalues.Value) (*apivalues.Value, error) {
				return apivalues.String("meow"), nil
			},
		),
	},
}
