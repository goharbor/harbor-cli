package pipeline

import (
	"dagger/harbor-cli/internal/dagger"
)

type Pipeline struct {
	source      *dagger.Directory
	dag         *dagger.Client
	appVersion  string
	goVersion   string
	GithubToken *dagger.Secret
}

func InitPipeline(source *dagger.Directory, dag *dagger.Client, appVersion string, goVersion string) *Pipeline {
	return &Pipeline{
		source:     source,
		dag:        dag,
		goVersion:  goVersion,
		appVersion: appVersion,
	}
}
