package version

var (
	Version   = "0.1.0"
	GitCommit = ""
)

func GetVersion() string {
	return Version
}
