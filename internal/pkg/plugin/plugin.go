package plugin

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"go.autokitteh.dev/sdk/api/apivalues"
	"go.autokitteh.dev/sdk/pluginimpl"
)

// TODO: This should not be a singleton. The caller to the plugin must
//       specify the configuration. This is for demo purposes only.
var awsConfig aws.Config

func init() {
	var err error
	if awsConfig, err = config.LoadDefaultConfig(context.Background()); err != nil {
		panic(err)
	}
}

var (
	NewPlugin = func() *pluginimpl.Plugin {
		return &pluginimpl.Plugin{
			ID:  "aws",
			Doc: "AWS SDK Plugin",
			Members: map[string]*pluginimpl.PluginMember{
				"ec2":                importService("EC2", ec2.NewFromConfig),
				"secret_credentials": credentials,
			},
		}
	}

	credentials = pluginimpl.NewSimpleMethodMember(
		"retreive AWS credentials",
		func(
			ctx context.Context,
			args []*apivalues.Value,
			kwargs map[string]*apivalues.Value,
		) (*apivalues.Value, error) {
			if err := pluginimpl.UnpackArgs(args, kwargs); err != nil {
				return nil, err
			}

			creds, err := awsConfig.Credentials.Retrieve(ctx)
			if err != nil {
				return nil, fmt.Errorf("retreive: %w", err)
			}

			return apivalues.Struct(
				apivalues.Symbol("aws_credentials"),
				map[string]*apivalues.Value{
					"AccessKeyID":     apivalues.String(creds.AccessKeyID),
					"SecretAccessKey": apivalues.String(creds.SecretAccessKey),
					"SessionToken":    apivalues.String(creds.SessionToken),
					"Source":          apivalues.String(creds.Source),
					"CanExpire":       apivalues.Boolean(creds.CanExpire),
					"Expires":         apivalues.Time(creds.Expires),
				},
			), nil
		},
	)
)
