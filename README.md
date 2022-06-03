# protoc-gen-go-errors

This protobuf plugins enables rpc method to list possible errors returned. The
plugin will generate single, package level, Go error variable for every error
listed. Server may use the generated error variable as return error and client
Can catch the error with `errors.Is()` assertion.

Every error definition **MUST** contain two properties:

1. arbitrary error code string
   This can be served as human friendly error identifier

2. an enumerated gRPC status [code](https://www.grpc.io/docs/guides/error/#error-status-codes)
   This will be transalated into real gRPC status code which will be used over the wire.

An simple example will be like below

```protobuf
syntax="proto3";

import "protoc-gen-go-errors/proto/error.proto"

// ...

service Sample {
  rpc Method(Input) returns (Response) {
    option(errors.types) = {
       status: ALREADY_EXISTS,
       code: "err_email_exists"
    };
    option(errors.types) = {
       status: INVALID_ARGUMENT,
       code: "err_invalid_email"
    };
  }
}
```

## Installation

1. Install protoc-gen-go-errors

```sh
$ go install github.com/ncrypthic/proto-gen-go-errors
```

2. Generate error definition Golang file

```sh
protoc --go-errors_out=:<dir> -I <path to protoc-gen-go-errors parent dir> -I <path to proto file directory> <proto_file_name>
```

## LICENSE

GNU General Public License v3.0 or later

See [COPYING](COPYING) to see the full text.
