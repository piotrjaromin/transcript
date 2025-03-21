package version

// Version information set by build flags
var (
	// Version is the current version of the application
	Version = "dev"
	
	// CommitSHA is the git commit SHA used to build the application
	CommitSHA = "unknown"
)

// GetVersion returns the current version information
func GetVersion() string {
	return Version
}

// GetCommitSHA returns the git commit SHA used to build the application
func GetCommitSHA() string {
	return CommitSHA
}
