//go:build unit

package plugin

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/assert"

	"go.autokitteh.dev/sdk/api/apivalues"
)

func TestImportService(t *testing.T) {
	ms := importServiceMethods(&ec2.Client{})
	assert.NotNil(t, ms)

	m, ok := ms["DescribeVpcs"]
	if !assert.True(t, ok) {
		return
	}

	v, err := m(
		context.Background(),
		"DescribeVpcs",
		[]*apivalues.Value{testRequestValue},
		nil,
		nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, v)
}
