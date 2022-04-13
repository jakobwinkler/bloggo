package util

// to set this, use
// `go build -ldflags "-X util.Version=$(git describe --always --tags --dirty)"`
var Version string = "(development version)"
