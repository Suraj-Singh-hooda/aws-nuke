package awsutil

// ErrSkipRequest type
type ErrSkipRequest string

// Error function
func (err ErrSkipRequest) Error() string {
	return string(err)
}

// ErrUnknownEndpoint type
type ErrUnknownEndpoint string

// Error function
func (err ErrUnknownEndpoint) Error() string {
	return string(err)
}
