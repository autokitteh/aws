module github.com/autokitteh/aws

go 1.18

// Uncomment these to build against local idl and sdk:
// replace go.autokitteh.dev/idl => ../idl
// replace go.autokitteh.dev/sdk => ../go-sdk
//
// RECOMMENDED: run ./scripts/git-hooks/install.sh to make sure these do not
// get comitted.

require (
	github.com/aws/aws-sdk-go v1.44.34 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
