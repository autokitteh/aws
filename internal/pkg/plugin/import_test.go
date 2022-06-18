//go:build real

package plugin

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"go.autokitteh.dev/sdk/api/apivalues"
)

var actualTestRequestValue = apivalues.DictFromMap(
	map[string]*apivalues.Value{
		"MaxResults": apivalues.Integer(100),
	},
)

func TestImportService(t *testing.T) {
	ms := importServiceMethods(ec2.NewFromConfig)
	assert.NotNil(t, ms)

	m, ok := ms["DescribeVpcs"]
	if !assert.True(t, ok) {
		return
	}

	v, err := m(
		context.Background(),
		"DescribeVpcs",
		[]*apivalues.Value{actualTestRequestValue},
		nil,
		nil,
	)
	if !assert.NoError(t, err) {
		return
	}

	assert.NotNil(t, v)

	// TODO: actually check if value is ok.
	spew.Dump(v)
}
