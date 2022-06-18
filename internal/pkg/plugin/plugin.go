package plugin

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

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
				"ec2": importService("EC2", &ec2.Client{}),
			},
		}
	}
)

func ec2Client() *ec2.Client { return ec2.NewFromConfig(awsConfig) }
