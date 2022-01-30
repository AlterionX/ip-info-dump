package base

import (
    "net"
    "errors"
)

type InfoResult struct {
    Info interface{}
    Err error
}

type InfoSource interface {
    // FetchInfo must close the returned channel if no messages will be sent through the channel.
    // If not, the system will hang forever.
    FetchInfo(arg net.IP) <-chan InfoResult
    Name() string
}

var BadArgument = errors.New("Requested IP/domain is malformed.")
var MissingKey = errors.New("No API key is available.")
var MissingReply = errors.New("No reply from source received.")
var SourceFailure = errors.New("Source failed to fetch information.")
