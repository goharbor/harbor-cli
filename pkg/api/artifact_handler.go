package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/scan"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// RunDeleteArtifact handles the deletion of an artifact.
func RunDeleteArtifact(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Artifact.DeleteArtifact(ctx, &artifact.DeleteArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to delete artifact: %v", err)
		return err
	}

	log.Infof("Artifact deleted successfully: %s/%s@%s", projectName, repoName, reference)
	return nil
}

// RunInfoArtifact retrieves information about a specific artifact.
func RunInfoArtifact(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.Artifact.GetArtifact(ctx, &artifact.GetArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to get artifact info: %v", err)
		return err
	}

	utils.PrintPayloadInJSONFormat(response.Payload)
	return nil
}

// RunListArtifact lists all artifacts in a repository.
func RunListArtifact(projectName, repoName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.Artifact.ListArtifacts(ctx, &artifact.ListArtifactsParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
	})
	if err != nil {
		log.Errorf("Failed to list artifacts: %v", err)
		return err
	}

	log.Info(response.Payload)
	return nil
}

// StartScanArtifact initiates a scan on a specific artifact.
func StartScanArtifact(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Scan.ScanArtifact(ctx, &scan.ScanArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to start scan: %v", err)
		return err
	}

	log.Infof("Scan started successfully: %s/%s@%s", projectName, repoName, reference)
	return nil
}

// StopScanArtifact stops a scan on a specific artifact.
func StopScanArtifact(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Scan.StopScanArtifact(ctx, &scan.StopScanArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to stop scan: %v", err)
		return err
	}

	log.Infof("Scan stopped successfully: %s/%s@%s", projectName, repoName, reference)
	return nil
}

// DeleteTag deletes a tag from a specific artifact.
func DeleteTag(projectName, repoName, reference, tag string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Artifact.DeleteTag(ctx, &artifact.DeleteTagParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
		TagName:        tag,
	})
	if err != nil {
		log.Errorf("Failed to delete tag: %v", err)
		return err
	}

	log.Infof("Tag deleted successfully: %s/%s@%s:%s", projectName, repoName, reference, tag)
	return nil
}

// ListTags lists all tags of a specific artifact.
func ListTags(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	resp, err := client.Artifact.ListTags(ctx, &artifact.ListTagsParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to list tags: %v", err)
		return err
	}

	log.Info(resp.Payload)
	return nil
}
