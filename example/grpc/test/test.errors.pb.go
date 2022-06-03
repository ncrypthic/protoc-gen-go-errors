package test

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var (
	ERR_err_email_exists  error = status.New(codes.AlreadyExists, "err_email_exists").Err()
	ERR_err_invalid_email error = status.New(codes.InvalidArgument, "err_invalid_email").Err()
)
