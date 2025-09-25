package pipeline

import (
	"dagger/harbor-cli/internal/dagger"
)

type Pipeline struct {
	source     *dagger.Directory
	dag        *dagger.Client
	appVersion string
}

func InitPipeline(source *dagger.Directory, dag *dagger.Client, appVersion string) *Pipeline {
	return &Pipeline{
		source:     source,
		dag:        dag,
		appVersion: appVersion,
	}
}
