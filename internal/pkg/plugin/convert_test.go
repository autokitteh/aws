//go:build unit

package plugin

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"

	"go.autokitteh.dev/sdk/api/apivalues"
)

var testRequestValue = apivalues.DictFromMap(
	map[string]*apivalues.Value{
		"dry_run":     apivalues.Boolean(false),
		"next_token":  apivalues.String("meow"),
		"max_results": apivalues.Integer(42),
		"vpc_ids":     apivalues.List(apivalues.String("first"), apivalues.String("second")),
		"filters": apivalues.List(
			apivalues.DictFromMap(map[string]*apivalues.Value{
				"name":   apivalues.String("gizmo"),
				"values": apivalues.List(apivalues.String("woof")),
			}),
			apivalues.DictFromMap(map[string]*apivalues.Value{
				"name": apivalues.String("zumi"),
			}),
		),
	},
)

var testRequest = ec2.DescribeVpcsInput{
	DryRun:     aws.Bool(false),
	MaxResults: aws.Int32(42),
	NextToken:  aws.String("meow"),
	VpcIds:     []string{"first", "second"},
	Filters: []types.Filter{
		types.Filter{
			Name:   aws.String("gizmo"),
			Values: []string{"woof"},
		},
		types.Filter{
			Name: aws.String("zumi"),
		},
	},
}

/*
func TestConvertFromAWS(t *testing.T) {
	v, err := ConvertFromAWS(nil, reflect.ValueOf(testRequest))
	if !assert.NoError(t, err) {
		return
	}

	// TODO: sort.
	if !assert.Equal(t, testRequestValue, v) {
		spew.Dump(testRequestValue)
		spew.Dump(v)
	}
}
*/

func TestConvertToAWS(t *testing.T) {
	var dvi ec2.DescribeVpcsInput

	if !assert.NoError(t, ConvertToAWS(nil, reflect.ValueOf(&dvi), testRequestValue)) {
		return
	}

	assert.Equal(t, testRequest, dvi)
}

// TODO: TestConvetFromAWS
