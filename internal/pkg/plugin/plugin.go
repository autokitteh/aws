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

var credentials *apivalues.Value

func init() {
	var err error
	if awsConfig, err = config.LoadDefaultConfig(context.Background()); err != nil {
		panic(err)
	}

	creds, err := awsConfig.Credentials.Retrieve(context.Background())
	if err != nil {
		panic(fmt.Errorf("credentials retreive: %w", err))
	}

	credentials = apivalues.Struct(
		apivalues.Symbol("aws_credentials"),
		map[string]*apivalues.Value{
			"access_key_id":     apivalues.String(creds.AccessKeyID),
			"secret_access_key": apivalues.String(creds.SecretAccessKey),
			"session_token":     apivalues.String(creds.SessionToken),
			"source":            apivalues.String(creds.Source),
			"can_expire":        apivalues.Boolean(creds.CanExpire),
			"expires":           apivalues.Time(creds.Expires),
		},
	)
}

var (
	NewPlugin = func() *pluginimpl.Plugin {
		return &pluginimpl.Plugin{
			ID:  "aws",
			Doc: "AWS SDK Plugin",
			Members: map[string]*pluginimpl.PluginMember{
				"ec2":                importService("EC2", ec2.NewFromConfig),
				"secret_credentials": pluginimpl.NewValueMember("AWS Credentials", credentials),
			},
		}
	}
)
