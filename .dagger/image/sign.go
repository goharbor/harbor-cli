package image

import (
	"context"
	"fmt"
)

func (s *ImagePipeline) Sign(ctx context.Context, imageAddr string) (string, error) {
	registryPasswordPlain, _ := s.RegistryPassword.Plaintext(ctx)

	cosing_ctr := s.dag.Container().From("cgr.dev/chainguard/cosign")

	// If githubToken is provided, use it to sign the image
	if s.GithubToken != nil {
		if s.ActionsIDTokenRequestURL == nil || s.ActionsIDTokenRequestToken == nil {
			return "", fmt.Errorf("actionsIdTokenRequestUrl (exist=%s) and actionsIdTokenRequestToken (exist=%t) must be provided when githubToken is provided",
				s.ActionsIDTokenRequestURL, s.ActionsIDTokenRequestToken != nil)
		}
		fmt.Printf("Setting the ENV Vars GITHUB_TOKEN, ACTIONS_ID_TOKEN_REQUEST_URL, ACTIONS_ID_TOKEN_REQUEST_TOKEN to sign with GitHub Token")
		cosing_ctr = cosing_ctr.WithSecretVariable("GITHUB_TOKEN", s.GithubToken).
			WithSecretVariable("ACTIONS_ID_TOKEN_REQUEST_URL", s.ActionsIDTokenRequestURL).
			WithSecretVariable("ACTIONS_ID_TOKEN_REQUEST_TOKEN", s.ActionsIDTokenRequestToken)
	}

	return cosing_ctr.WithSecretVariable("REGISTRY_PASSWORD", s.RegistryPassword).
		WithExec([]string{"cosign", "env"}).
		WithExec([]string{
			"cosign", "sign", "--yes", "--recursive",
			"--registry-username", s.RegistryUsername,
			"--registry-password", registryPasswordPlain,
			imageAddr,
			"--timeout", "1m",
		}).Stdout(ctx)
}
