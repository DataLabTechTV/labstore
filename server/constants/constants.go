package constants

// general constants
const (
	Empty   = "**EMPTY**"
	Unknown = "unknown"
)

// app constants
const (
	Name        = "LabStore"
	Description = "An S3-Compatible Object Store"
	Author      = "IllumiKnow Labs"
	Version     = "0.1.0"
)

// git repo constants
const (
	GitRepo           = "https://github.com/IllumiKnowLabs/labstore"
	GitAssetsFilename = "assets.zip"
)

// ldflags variables (defaults)
var (
	GitTag    = Unknown
	GitCommit = Unknown
	BuildTime = Unknown
	Builder   = Unknown
)
