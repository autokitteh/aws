#!/bin/bash

set -euo pipefail

AK_PLUGIN_ID="autokitteh.aws" AWSPLUGIND_GRPC_PORT=30001 AWSPLUGIND_HTTP_PORT=30000 exec bin/awsplugind "$@"
