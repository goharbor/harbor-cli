package image

import "dagger/harbor-cli/internal/dagger"

type ImagePipeline struct {
	source                     *dagger.Directory
	dag                        *dagger.Client
	RegistryPassword           *dagger.Secret
	GithubToken                *dagger.Secret
	ActionsIDTokenRequestURL   *dagger.Secret
	ActionsIDTokenRequestToken *dagger.Secret
	RegistryAddress            string
	RegistryUsername           string
	appVersion                 string
	goVersion                  string
}

func InitImagePipeline(dag *dagger.Client, source *dagger.Directory,
	registryPassword, githubToken, actionsIDTokenRequestURL, actionsIDTokenRequestToken *dagger.Secret,
	registryAddress, registryUsername, appVersion, goVersion string,
) *ImagePipeline {
	return &ImagePipeline{
		source:                     source,
		dag:                        dag,
		RegistryPassword:           registryPassword,
		GithubToken:                githubToken,
		RegistryAddress:            registryAddress,
		RegistryUsername:           registryUsername,
		ActionsIDTokenRequestURL:   actionsIDTokenRequestURL,
		ActionsIDTokenRequestToken: actionsIDTokenRequestToken,
		appVersion:                 appVersion,
		goVersion:                  goVersion,
	}
}
